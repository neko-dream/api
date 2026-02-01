package talksession_usecase_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/application/usecase/talksession_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/analysis"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/di"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/internal/test/txtest"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

// AnalysisRepositoryのモック
type analysisRepositoryMock struct{}

func (m *analysisRepositoryMock) FindByTalkSessionID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*analysis.AnalysisReport, error) {
	return nil, sql.ErrNoRows
}

func (m *analysisRepositoryMock) FindByID(ctx context.Context, analysisReportID shared.UUID[analysis.AnalysisReport]) (*analysis.AnalysisReport, error) {
	return nil, sql.ErrNoRows
}

func (m *analysisRepositoryMock) SaveReport(ctx context.Context, report *analysis.AnalysisReport) error {
	return nil
}

func newAnalysisRepositoryMock() analysis.AnalysisRepository {
	return &analysisRepositoryMock{}
}

type ToggleAnalysisUseCaseTestData struct {
	ToggleAnalysisUseCase talksession_usecase.ToggleAnalysisUseCase
	TalkSessionRepo       talksession.TalkSessionRepository
	UserRepo              user.UserRepository
	OwnerUser             *user.User
	NonOwnerUser          *user.User
	TalkSession           *talksession.TalkSession
	TalkSessionID         shared.UUID[talksession.TalkSession]
	OwnerUserID           shared.UUID[user.User]
	NonOwnerUserID        shared.UUID[user.User]
}

func TestToggleAnalysisUseCase_Execute(t *testing.T) {
	// 必要な環境変数を設定
	t.Setenv("ENCRYPTION_VERSION", "v1")
	t.Setenv("ENCRYPTION_SECRET", "12345678901234567890123456789012")
	t.Setenv("DATABASE_URL", "postgres://kotohiro:kotohiro@localhost:5432/kotohiro?sslmode=disable")
	t.Setenv("ENV", "test")

	// OpenTelemetryのnoop tracerを設定
	otel.SetTracerProvider(noop.NewTracerProvider())

	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)
	encryptor, _ := crypto.NewEncryptor(lo.ToPtr(config.Config{
		ENCRYPTION_VERSION: crypto.Version1,
		ENCRYPTION_SECRET:  "12345678901234567890123456789012", // テスト用の32バイトキー
	}))

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(
		dbManager,
		repository.NewImageRepositoryMock(),
		encryptor,
	)
	eventStore := repository.NewEventStoreMock()
	talkSessionRepo := repository.NewTalkSessionRepository(dbManager, eventStore)

	// テスト用のconfig作成（LOCAL環境に設定）
	testConfig := &config.Config{
		Env: config.LOCAL, // ローカル環境として設定
	}

	// ユースケースの初期化
	toggleUseCase := talksession_usecase.NewToggleAnalysisUseCase(
		talkSessionRepo,
		repository.NewAnalysisServiceMock(),
		newAnalysisRepositoryMock(),
		dbManager,
		testConfig,
	)

	// テストデータの初期化
	initData := &ToggleAnalysisUseCaseTestData{
		ToggleAnalysisUseCase: toggleUseCase,
		TalkSessionRepo:       talkSessionRepo,
		UserRepo:              userRepo,
		TalkSessionID:         shared.NewUUID[talksession.TalkSession](),
		OwnerUserID:           shared.NewUUID[user.User](),
		NonOwnerUserID:        shared.NewUUID[user.User](),
	}

	testCases := []*txtest.TransactionalTestCase[ToggleAnalysisUseCaseTestData]{
		{
			Name: "正常系: 分析無効化（false→true）",
			SetupFn: func(ctx context.Context, data *ToggleAnalysisUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				// トークセッション作成（disableAnalysis=false）
				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"テストテーマ",
					lo.ToPtr("テスト説明"),
					nil, // thumbnailURL
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,                // location
					lo.ToPtr("渋谷区"), // city
					lo.ToPtr("東京都"), // prefecture
					false,              // hideTop
					nil,                // organizationID
					nil,                // organizationAliasID
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *ToggleAnalysisUseCaseTestData) error {
				input := talksession_usecase.ToggleAnalysisInput{
					TalkSessionID: data.TalkSessionID,
					UserID:        data.OwnerUserID,
				}

				output, err := data.ToggleAnalysisUseCase.Execute(ctx, input)
				if err != nil {
					return err
				}

				// false → true にトグル
				assert.True(t, output.DisableAnalysis)

				// DBから再取得して確認
				updated, err := data.TalkSessionRepo.FindByID(ctx, data.TalkSessionID)
				if err != nil {
					return err
				}
				assert.True(t, updated.DisableAnalysis())

				return nil
			},
			WantErr: false,
		},
		{
			Name: "正常系: 分析有効化（true→false）",
			SetupFn: func(ctx context.Context, data *ToggleAnalysisUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				// トークセッション作成（disableAnalysis=true）
				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"テストテーマ",
					lo.ToPtr("テスト説明"),
					nil, // thumbnailURL
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,                // location
					lo.ToPtr("渋谷区"), // city
					lo.ToPtr("東京都"), // prefecture
					false,              // hideTop
					nil,                // organizationID
					nil,                // organizationAliasID
				)
				data.TalkSession.SetDisableAnalysis(true) // 分析無効化状態にする
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *ToggleAnalysisUseCaseTestData) error {
				input := talksession_usecase.ToggleAnalysisInput{
					TalkSessionID: data.TalkSessionID,
					UserID:        data.OwnerUserID,
				}

				output, err := data.ToggleAnalysisUseCase.Execute(ctx, input)
				if err != nil {
					return err
				}

				// true → false にトグル
				assert.False(t, output.DisableAnalysis)

				// DBから再取得して確認
				updated, err := data.TalkSessionRepo.FindByID(ctx, data.TalkSessionID)
				if err != nil {
					return err
				}
				assert.False(t, updated.DisableAnalysis())

				return nil
			},
			WantErr: false,
		},
		{
			Name: "正常系: 2回トグルすると元に戻る",
			SetupFn: func(ctx context.Context, data *ToggleAnalysisUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"テストテーマ",
					lo.ToPtr("テスト説明"),
					nil, // thumbnailURL
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,                // location
					lo.ToPtr("渋谷区"), // city
					lo.ToPtr("東京都"), // prefecture
					false,              // hideTop
					nil,                // organizationID
					nil,                // organizationAliasID
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *ToggleAnalysisUseCaseTestData) error {
				input := talksession_usecase.ToggleAnalysisInput{
					TalkSessionID: data.TalkSessionID,
					UserID:        data.OwnerUserID,
				}

				// 1回目のトグル（false → true）
				output1, err := data.ToggleAnalysisUseCase.Execute(ctx, input)
				if err != nil {
					return err
				}
				assert.True(t, output1.DisableAnalysis)

				// 2回目のトグル（true → false）
				output2, err := data.ToggleAnalysisUseCase.Execute(ctx, input)
				if err != nil {
					return err
				}
				assert.False(t, output2.DisableAnalysis)

				return nil
			},
			WantErr: false,
		},
		{
			Name: "異常系: 存在しないトークセッション",
			TestFn: func(ctx context.Context, data *ToggleAnalysisUseCaseTestData) error {
				nonExistentID := shared.NewUUID[talksession.TalkSession]()
				input := talksession_usecase.ToggleAnalysisInput{
					TalkSessionID: nonExistentID,
					UserID:        data.OwnerUserID,
				}

				_, err := data.ToggleAnalysisUseCase.Execute(ctx, input)
				assert.ErrorIs(t, err, messages.TalkSessionNotFound)
				return nil
			},
			WantErr: false,
		},
		{
			Name: "正常系: ローカル環境では非オーナーでもトグル可能",
			SetupFn: func(ctx context.Context, data *ToggleAnalysisUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"テストテーマ",
					lo.ToPtr("テスト説明"),
					nil, // thumbnailURL
					data.OwnerUserID, // オーナーは別ユーザー
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,                // location
					lo.ToPtr("渋谷区"), // city
					lo.ToPtr("東京都"), // prefecture
					false,              // hideTop
					nil,                // organizationID
					nil,                // organizationAliasID
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *ToggleAnalysisUseCaseTestData) error {
				input := talksession_usecase.ToggleAnalysisInput{
					TalkSessionID: data.TalkSessionID,
					UserID:        data.NonOwnerUserID, // 非オーナー
				}

				// ローカル環境では成功する
				output, err := data.ToggleAnalysisUseCase.Execute(ctx, input)
				assert.NoError(t, err)
				assert.True(t, output.DisableAnalysis)
				return nil
			},
			WantErr: false,
		},
	}

	// テストケースの実行
	txtest.RunTransactionalTests(t, dbManager, initData, testCases)
}

// プロダクション環境でのテスト
func TestToggleAnalysisUseCase_ProductionEnvironment(t *testing.T) {
	// プロダクション用のconfig作成
	prodConfig := &config.Config{
		Env: config.PROD, // プロダクション環境として設定
	}

	// 環境変数設定
	t.Setenv("ENCRYPTION_VERSION", "v1")
	t.Setenv("ENCRYPTION_SECRET", "12345678901234567890123456789012")
	t.Setenv("DATABASE_URL", "postgres://kotohiro:kotohiro@localhost:5432/kotohiro?sslmode=disable")
	t.Setenv("ENV", "production")

	// OpenTelemetryのnoop tracerを設定
	otel.SetTracerProvider(noop.NewTracerProvider())

	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)
	encryptor, _ := crypto.NewEncryptor(lo.ToPtr(config.Config{
		ENCRYPTION_VERSION: crypto.Version1,
		ENCRYPTION_SECRET:  "12345678901234567890123456789012",
	}))

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(
		dbManager,
		repository.NewImageRepositoryMock(),
		encryptor,
	)
	eventStore := repository.NewEventStoreMock()
	talkSessionRepo := repository.NewTalkSessionRepository(dbManager, eventStore)

	// プロダクション環境用のユースケースの初期化
	toggleUseCase := talksession_usecase.NewToggleAnalysisUseCase(
		talkSessionRepo,
		repository.NewAnalysisServiceMock(),
		newAnalysisRepositoryMock(),
		dbManager,
		prodConfig,
	)

	err := dbManager.TestTx(context.Background(), func(ctx context.Context) error {
		// テストデータ準備
		ownerUserID := shared.NewUUID[user.User]()
		nonOwnerUserID := shared.NewUUID[user.User]()
		talkSessionID := shared.NewUUID[talksession.TalkSession]()

		// オーナーユーザーの作成
		ownerUser := user.NewUser(
			ownerUserID,
			lo.ToPtr("owner_display_id"),
			lo.ToPtr("Owner User"),
			"owner@example.com",
			"GOOGLE",
			lo.ToPtr("https://example.com/owner-icon.jpg"),
		)
		if err := userRepo.Create(ctx, ownerUser); err != nil {
			return err
		}

		// トークセッション作成
		talkSession := talksession.NewTalkSession(
			talkSessionID,
			"テストテーマ",
			lo.ToPtr("テスト説明"),
			nil, // thumbnailURL
			ownerUserID,
			clock.Now(ctx),
			clock.Now(ctx).Add(time.Hour*24),
			nil,                // location
			lo.ToPtr("渋谷区"), // city
			lo.ToPtr("東京都"), // prefecture
			false,              // hideTop
			nil,                // organizationID
			nil,                // organizationAliasID
		)
		if err := talkSessionRepo.Create(ctx, talkSession); err != nil {
			return err
		}

		// プロダクション環境では非オーナーがトグルを試みると失敗する
		input := talksession_usecase.ToggleAnalysisInput{
			TalkSessionID: talkSessionID,
			UserID:        nonOwnerUserID, // 非オーナーのID
		}

		_, err := toggleUseCase.Execute(ctx, input)
		assert.ErrorIs(t, err, messages.ForbiddenError)

		return nil
	})

	assert.NoError(t, err)
}
