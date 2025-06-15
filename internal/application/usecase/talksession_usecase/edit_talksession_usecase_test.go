package talksession_usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/neko-dream/server/internal/application/usecase/talksession_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/crypto"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/server/internal/test/txtest"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

type EditTalkSessionUseCaseTestData struct {
	EditTalkSessionUseCase talksession_usecase.EditTalkSessionUseCase
	TalkSessionRepo        talksession.TalkSessionRepository
	UserRepo               user.UserRepository
	OwnerUser              *user.User
	NonOwnerUser           *user.User
	TalkSession            *talksession.TalkSession
	TalkSessionID          shared.UUID[talksession.TalkSession]
	OwnerUserID            shared.UUID[user.User]
	NonOwnerUserID         shared.UUID[user.User]
}

func TestEditTalkSessionUseCase_Execute(t *testing.T) {
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
	talkSessionRepo := repository.NewTalkSessionRepository(dbManager)

	// テスト用のconfig作成（LOCAL環境に設定）
	testConfig := &config.Config{
		Env: config.LOCAL, // ローカル環境として設定
	}

	// コマンドハンドラの初期化
	editCommand := talksession_usecase.NewEditTalkSessionUseCase(
		talkSessionRepo,
		userRepo,
		dbManager,
		testConfig,
	)

	// テストデータの初期化
	initData := &EditTalkSessionUseCaseTestData{
		EditTalkSessionUseCase: editCommand,
		TalkSessionRepo:        talkSessionRepo,
		UserRepo:               userRepo,
		TalkSessionID:          shared.NewUUID[talksession.TalkSession](),
		OwnerUserID:            shared.NewUUID[user.User](),
		NonOwnerUserID:         shared.NewUUID[user.User](),
	}

	testCases := []*txtest.TransactionalTestCase[EditTalkSessionUseCaseTestData]{
		{
			Name: "正常系: オーナーがトークセッションを編集できる",
			SetupFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				data.OwnerUser = &ownerUser
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				// 初期のトークセッション作成
				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"初期テーマ",
					lo.ToPtr("初期説明文"),
					lo.ToPtr("https://example.com/initial.jpg"),
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil, // location
					lo.ToPtr("初期市区町村"),
					lo.ToPtr("初期都道府県"),
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// 編集用の入力データ
				input := talksession_usecase.EditTalkSessionInput{
					TalkSessionID:    data.TalkSessionID,
					UserID:           data.OwnerUserID,
					Theme:            "編集後のテーマ",
					Description:      lo.ToPtr("編集後の説明文"),
					ThumbnailURL:     lo.ToPtr("https://example.com/edited.jpg"),
					ScheduledEndTime: clock.Now(ctx).Add(time.Hour * 48),
					Latitude:         lo.ToPtr(35.6895),
					Longitude:        lo.ToPtr(139.6917),
					City:             lo.ToPtr("渋谷区"),
					Prefecture:       lo.ToPtr("東京都"),
				}

				// コマンドの実行
				output, err := data.EditTalkSessionUseCase.Execute(ctx, input)
				if err != nil {
					return err
				}

				// 結果の検証
				assert.Equal(t, input.Theme, output.TalkSession.Theme)
				assert.Equal(t, input.Description, output.TalkSession.Description)
				assert.Equal(t, input.ThumbnailURL, output.TalkSession.ThumbnailURL)
				assert.Equal(t, input.ScheduledEndTime, output.TalkSession.ScheduledEndTime)
				assert.Equal(t, input.City, output.TalkSession.City)
				assert.Equal(t, input.Prefecture, output.TalkSession.Prefecture)
				assert.Equal(t, input.Latitude, output.Latitude)
				assert.Equal(t, input.Longitude, output.Longitude)
				// 制限事項は編集では更新されないので、元のまま
				assert.Empty(t, output.Restrictions)

				// DBから再取得して確認
				updatedTalkSession, err := data.TalkSessionRepo.FindByID(ctx, data.TalkSessionID)
				if err != nil {
					return err
				}
				assert.Equal(t, input.Theme, updatedTalkSession.Theme())
				assert.Equal(t, input.Description, updatedTalkSession.Description())
				assert.Equal(t, input.City, updatedTalkSession.City())
				assert.Equal(t, input.Prefecture, updatedTalkSession.Prefecture())

				return nil
			},
			WantErr: false,
		},
		{
			Name: "正常系: ローカル環境では非オーナーでも編集できる",
			SetupFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				data.OwnerUser = &ownerUser
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				// 非オーナーユーザーの作成
				nonOwnerUser := user.NewUser(
					data.NonOwnerUserID,
					lo.ToPtr("nonowner"),
					lo.ToPtr("Non-Owner User"),
					"nonowner@example.com",
					"GOOGLE",
					nil,
				)
				data.NonOwnerUser = &nonOwnerUser
				if err := data.UserRepo.Create(ctx, nonOwnerUser); err != nil {
					return err
				}

				// トークセッション作成（オーナーが作成）
				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"初期テーマ",
					lo.ToPtr("初期説明文"),
					lo.ToPtr("https://example.com/initial.jpg"),
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,
					lo.ToPtr("初期市区町村"),
					lo.ToPtr("初期都道府県"),
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// 非オーナーが編集を試みる
				input := talksession_usecase.EditTalkSessionInput{
					TalkSessionID:    data.TalkSessionID,
					UserID:           data.NonOwnerUserID, // 非オーナーのID
					Theme:            "不正な編集",
					ScheduledEndTime: clock.Now(ctx).Add(time.Hour * 24),
				}

				_, err := data.EditTalkSessionUseCase.Execute(ctx, input)
				assert.NoError(t, err) // ローカル環境では成功する
				return nil
			},
			WantErr: false,
		},
		{
			Name: "異常系: 存在しないトークセッションは編集できない",
			TestFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				nonExistentID := shared.NewUUID[talksession.TalkSession]()
				input := talksession_usecase.EditTalkSessionInput{
					TalkSessionID:    nonExistentID,
					UserID:           data.OwnerUserID,
					Theme:            "編集",
					ScheduledEndTime: clock.Now(ctx).Add(time.Hour * 24),
				}

				_, err := data.EditTalkSessionUseCase.Execute(ctx, input)
				assert.ErrorIs(t, err, messages.TalkSessionNotFound)
				return nil
			},
			WantErr: false,
		},
		{
			Name: "異常系: テーマが100文字を超える場合はエラー",
			SetupFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				data.OwnerUser = &ownerUser
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				// 初期のトークセッション作成
				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"初期テーマ",
					lo.ToPtr("初期説明文"),
					lo.ToPtr("https://example.com/initial.jpg"),
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,
					lo.ToPtr("初期市区町村"),
					lo.ToPtr("初期都道府県"),
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// 101文字のテーマを生成
				longTheme := ""
				for i := 0; i < 101; i++ {
					longTheme += "あ"
				}

				input := talksession_usecase.EditTalkSessionInput{
					TalkSessionID:    data.TalkSessionID,
					UserID:           data.OwnerUserID,
					Theme:            longTheme,
					ScheduledEndTime: clock.Now(ctx).Add(time.Hour * 24),
				}

				_, err := data.EditTalkSessionUseCase.Execute(ctx, input)
				assert.ErrorIs(t, err, messages.TalkSessionThemeTooLong)
				return nil
			},
			WantErr: false,
		},
		{
			Name: "異常系: 説明文が40000文字を超える場合はエラー",
			SetupFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				data.OwnerUser = &ownerUser
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				// 初期のトークセッション作成
				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"初期テーマ",
					lo.ToPtr("初期説明文"),
					lo.ToPtr("https://example.com/initial.jpg"),
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,
					lo.ToPtr("初期市区町村"),
					lo.ToPtr("初期都道府県"),
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// 40001文字の説明文を生成
				longDescription := ""
				for i := 0; i < 40001; i++ {
					longDescription += "あ"
				}

				input := talksession_usecase.EditTalkSessionInput{
					TalkSessionID:    data.TalkSessionID,
					UserID:           data.OwnerUserID,
					Theme:            "テーマ",
					Description:      &longDescription,
					ScheduledEndTime: clock.Now(ctx).Add(time.Hour * 24),
				}

				_, err := data.EditTalkSessionUseCase.Execute(ctx, input)
				assert.ErrorIs(t, err, messages.TalkSessionDescriptionTooLong)
				return nil
			},
			WantErr: false,
		},
		{
			Name: "異常系: 終了時刻が過去の場合はエラー",
			SetupFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				data.OwnerUser = &ownerUser
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				// 初期のトークセッション作成
				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"初期テーマ",
					lo.ToPtr("初期説明文"),
					lo.ToPtr("https://example.com/initial.jpg"),
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,
					lo.ToPtr("初期市区町村"),
					lo.ToPtr("初期都道府県"),
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				input := talksession_usecase.EditTalkSessionInput{
					TalkSessionID:    data.TalkSessionID,
					UserID:           data.OwnerUserID,
					Theme:            "テーマ",
					ScheduledEndTime: clock.Now(ctx).Add(-time.Hour), // 過去の時刻
				}

				_, err := data.EditTalkSessionUseCase.Execute(ctx, input)
				assert.ErrorIs(t, err, messages.InvalidScheduledEndTime)
				return nil
			},
			WantErr: false,
		},
		{
			Name: "正常系: 位置情報なしでも編集できる",
			SetupFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				data.OwnerUser = &ownerUser
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				// 初期のトークセッション作成
				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"初期テーマ",
					lo.ToPtr("初期説明文"),
					lo.ToPtr("https://example.com/initial.jpg"),
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil, // location
					lo.ToPtr("初期市区町村"),
					lo.ToPtr("初期都道府県"),
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				input := talksession_usecase.EditTalkSessionInput{
					TalkSessionID:    data.TalkSessionID,
					UserID:           data.OwnerUserID,
					Theme:            "位置情報なしテーマ",
					ScheduledEndTime: clock.Now(ctx).Add(time.Hour * 24),
					// Latitude, Longitudeを指定しない
				}

				output, err := data.EditTalkSessionUseCase.Execute(ctx, input)
				if err != nil {
					return err
				}

				assert.Equal(t, input.Theme, output.TalkSession.Theme)
				assert.Nil(t, output.Latitude)
				assert.Nil(t, output.Longitude)
				return nil
			},
			WantErr: false,
		},
		{
			Name: "正常系: 部分的な更新も可能",
			SetupFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// オーナーユーザーの作成
				ownerUser := user.NewUser(
					data.OwnerUserID,
					lo.ToPtr("owner_display_id"),
					lo.ToPtr("Owner User"),
					"owner@example.com",
					"GOOGLE",
					lo.ToPtr("https://example.com/owner-icon.jpg"),
				)
				data.OwnerUser = &ownerUser
				if err := data.UserRepo.Create(ctx, ownerUser); err != nil {
					return err
				}

				// 初期のトークセッション作成
				data.TalkSession = talksession.NewTalkSession(
					data.TalkSessionID,
					"初期テーマ",
					lo.ToPtr("初期説明文"),
					lo.ToPtr("https://example.com/initial.jpg"),
					data.OwnerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil, // location
					lo.ToPtr("初期市区町村"),
					lo.ToPtr("初期都道府県"),
				)
				return data.TalkSessionRepo.Create(ctx, data.TalkSession)
			},
			TestFn: func(ctx context.Context, data *EditTalkSessionUseCaseTestData) error {
				// テーマと終了時刻のみ更新
				input := talksession_usecase.EditTalkSessionInput{
					TalkSessionID:    data.TalkSessionID,
					UserID:           data.OwnerUserID,
					Theme:            "部分更新後のテーマ",
					ScheduledEndTime: clock.Now(ctx).Add(time.Hour * 72),
					// 他のフィールドはnilのまま
				}

				output, err := data.EditTalkSessionUseCase.Execute(ctx, input)
				if err != nil {
					return err
				}

				// 更新されたフィールドの確認
				assert.Equal(t, input.Theme, output.TalkSession.Theme)
				assert.Equal(t, input.ScheduledEndTime, output.TalkSession.ScheduledEndTime)

				// EditTalkSessionUseCaseInputで指定していないフィールドはnilになる
				assert.Nil(t, output.TalkSession.Description)
				assert.Nil(t, output.TalkSession.ThumbnailURL)

				return nil
			},
			WantErr: false,
		},
	}

	// テストケースの実行
	txtest.RunTransactionalTests(t, dbManager, initData, testCases)
}

// TestEditTalkSessionInput_Validate 入力値バリデーションのユニットテスト
func TestEditTalkSessionInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   talksession_usecase.EditTalkSessionInput
		wantErr error
	}{
		{
			name: "正常系: 有効な入力",
			input: talksession_usecase.EditTalkSessionInput{
				Theme:       "正常なテーマ",
				Description: lo.ToPtr("正常な説明文"),
			},
			wantErr: nil,
		},
		{
			name: "異常系: テーマが100文字を超える",
			input: talksession_usecase.EditTalkSessionInput{
				Theme: func() string {
					s := ""
					for i := 0; i < 101; i++ {
						s += "あ"
					}
					return s
				}(),
			},
			wantErr: messages.TalkSessionThemeTooLong,
		},
		{
			name: "異常系: 説明文が40000文字を超える",
			input: talksession_usecase.EditTalkSessionInput{
				Theme: "正常なテーマ",
				Description: func() *string {
					s := ""
					for i := 0; i < 40001; i++ {
						s += "あ"
					}
					return &s
				}(),
			},
			wantErr: messages.TalkSessionDescriptionTooLong,
		},
		{
			name: "正常系: 説明文がnil",
			input: talksession_usecase.EditTalkSessionInput{
				Theme:       "正常なテーマ",
				Description: nil,
			},
			wantErr: nil,
		},
		{
			name: "正常系: 境界値 - テーマがちょうど100文字",
			input: talksession_usecase.EditTalkSessionInput{
				Theme: func() string {
					s := ""
					for i := 0; i < 100; i++ {
						s += "あ"
					}
					return s
				}(),
			},
			wantErr: nil,
		},
		{
			name: "正常系: 境界値 - 説明文がちょうど40000文字",
			input: talksession_usecase.EditTalkSessionInput{
				Theme: "正常なテーマ",
				Description: func() *string {
					s := ""
					for i := 0; i < 40000; i++ {
						s += "あ"
					}
					return &s
				}(),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// プロダクション環境でのテスト
func TestEditTalkSessionUseCase_ProductionEnvironment(t *testing.T) {
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
	talkSessionRepo := repository.NewTalkSessionRepository(dbManager)

	// プロダクション環境用のコマンドハンドラの初期化
	editCommand := talksession_usecase.NewEditTalkSessionUseCase(
		talkSessionRepo,
		userRepo,
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

		// 非オーナーユーザーの作成
		nonOwnerUser := user.NewUser(
			nonOwnerUserID,
			lo.ToPtr("nonowner"),
			lo.ToPtr("Non-Owner User"),
			"nonowner@example.com",
			"GOOGLE",
			nil,
		)
		if err := userRepo.Create(ctx, nonOwnerUser); err != nil {
			return err
		}

		// トークセッション作成
		talkSession := talksession.NewTalkSession(
			talkSessionID,
			"初期テーマ",
			lo.ToPtr("初期説明文"),
			lo.ToPtr("https://example.com/initial.jpg"),
			ownerUserID,
			clock.Now(ctx),
			clock.Now(ctx).Add(time.Hour*24),
			nil,
			lo.ToPtr("初期市区町村"),
			lo.ToPtr("初期都道府県"),
		)
		if err := talkSessionRepo.Create(ctx, talkSession); err != nil {
			return err
		}

		// プロダクション環境では非オーナーが編集を試みると失敗する
		input := talksession_usecase.EditTalkSessionInput{
			TalkSessionID:    talkSessionID,
			UserID:           nonOwnerUserID, // 非オーナーのID
			Theme:            "不正な編集",
			ScheduledEndTime: clock.Now(ctx).Add(time.Hour * 24),
		}

		_, err := editCommand.Execute(ctx, input)
		assert.ErrorIs(t, err, messages.ForbiddenError)

		return nil
	})

	assert.NoError(t, err)
}
