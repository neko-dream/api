package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/crypto"
	ci "github.com/neko-dream/server/internal/infrastructure/crypto"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/server/internal/test/txtest"
	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"
)

func TestTalkSessionRepository_Create(t *testing.T) {
	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)
	encryptor, _ := ci.NewEncryptor(lo.ToPtr(config.Config{
		ENCRYPTION_VERSION: crypto.Version1,
		ENCRYPTION_SECRET:  "12345678901234567890123456789012", // テスト用の32バイトキー
	}))
	type TestData struct {
		TsRepo      talksession.TalkSessionRepository
		TalkSession *talksession.TalkSession
	}

	initData := TestData{}
	talkSessionID := shared.NewUUID[talksession.TalkSession]()
	ownerUserID := shared.NewUUID[user.User]()
	userRepo := repository.NewUserRepository(
		dbManager,
		repository.NewImageRepositoryMock(),
		encryptor,
	)
	userRepo.Create(context.Background(), user.NewUser(
		ownerUserID,
		lo.ToPtr("test"),
		lo.ToPtr("test"),
		"test",
		"GOOGLE",
		nil,
	))

	testCases := []*txtest.TransactionalTestCase[TestData]{
		{
			Name: "トークセッション作成ができる",
			SetupFn: func(ctx context.Context, data *TestData) error {
				data.TalkSession = talksession.NewTalkSession(
					talkSessionID,
					"test",
					nil,
					lo.ToPtr("https://example.com/test.jpg"),
					ownerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,
					nil,
					nil,
				)
				data.TsRepo = repository.NewTalkSessionRepository(dbManager)
				return nil
			},
			TestFn: func(ctx context.Context, data *TestData) error {
				if err := data.TsRepo.Create(ctx, data.TalkSession); err != nil {
					return err
				}

				ts, err := dbManager.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
				if err != nil {
					return err
				}
				if ts.TalkSession.TalkSessionID != talkSessionID.UUID() {
					return errors.New("トークセッションIDが一致しません")
				}

				return nil
			},
			WantErr: false,
		},
		{
			Name: "トークセッション作成ができ、Locationも保存される",
			SetupFn: func(ctx context.Context, data *TestData) error {
				data.TalkSession = talksession.NewTalkSession(
					talkSessionID,
					"test",
					nil,
					lo.ToPtr("https://example.com/test.jpg"),
					ownerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					talksession.NewLocation(
						talkSessionID,
						30.0,
						30.0,
					),
					nil, nil,
				)
				data.TsRepo = repository.NewTalkSessionRepository(dbManager)
				return nil
			},
			TestFn: func(ctx context.Context, data *TestData) error {
				if err := data.TsRepo.Create(ctx, data.TalkSession); err != nil {
					return err
				}

				ts, err := dbManager.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
				if err != nil {
					return err
				}
				if ts.TalkSession.TalkSessionID != talkSessionID.UUID() {
					return errors.New("トークセッションIDが一致しません")
				}
				location := talksession.NewLocation(
					talkSessionID,
					ts.Latitude,
					ts.Longitude,
				)
				assert.Equal(t, location.ToGeographyText(), "POINT(30.000000 30.000000)")

				return nil
			},
			WantErr: false,
		},
	}

	txtest.RunTransactionalTests(t, dbManager, &initData, testCases)
}
