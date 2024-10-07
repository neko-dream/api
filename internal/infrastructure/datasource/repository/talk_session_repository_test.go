package repository_test

import (
	"errors"
	"log"
	"testing"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/datasource/repository"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/test/txtest"

	"github.com/stretchr/testify/assert"
)

func TestTalkSessionRepository_Create(t *testing.T) {
	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)

	type TestData struct {
		TsRepo      talksession.TalkSessionRepository
		TalkSession *talksession.TalkSession
	}

	initData := TestData{
		TsRepo: repository.NewTalkSessionRepository(dbManager),
	}
	talkSessionID := shared.NewUUID[talksession.TalkSession]()
	ownerUserID := shared.NewUUID[user.User]()
	testCases := []*txtest.TransactionalTestCase[TestData]{
		{
			Name: "トークセッション作成ができる",
			SetupFn: func(ctx *txtest.TestContext[TestData]) error {
				ctx.Data.TalkSession = talksession.NewTalkSession(
					talkSessionID,
					"test",
					ownerUserID,
					nil,
					time.Now(ctx),
					// 明日
					time.Now(ctx).Add(ctx, time.Hour*24),
					nil,
				)
				ctx.Data.TsRepo = repository.NewTalkSessionRepository(dbManager)
				return nil
			},
			TestFn: func(ctx *txtest.TestContext[TestData]) error {
				if err := ctx.Data.TsRepo.Create(ctx, ctx.Data.TalkSession); err != nil {
					return err
				}

				ts, err := dbManager.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
				if err != nil {
					return err
				}
				if ts.TalkSessionID != talkSessionID.UUID() {
					return errors.New("トークセッションIDが一致しません")
				}

				log.Println(ts)

				return nil
			},
			WantErr: false,
		},
		{
			Name: "トークセッション作成ができ、Locationも保存される",
			SetupFn: func(ctx *txtest.TestContext[TestData]) error {
				ctx.Data.TalkSession = talksession.NewTalkSession(
					talkSessionID,
					"test",
					ownerUserID,
					nil,
					time.Now(ctx),
					time.Now(ctx).Add(ctx, time.Hour*24),
					talksession.NewLocation(
						talkSessionID,
						30.0,
						30.0,
						"鯖江市",
						"福井県",
					),
				)
				ctx.Data.TsRepo = repository.NewTalkSessionRepository(dbManager)
				return nil
			},
			TestFn: func(ctx *txtest.TestContext[TestData]) error {
				if err := ctx.Data.TsRepo.Create(ctx, ctx.Data.TalkSession); err != nil {
					return err
				}

				ts, err := dbManager.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
				if err != nil {
					return err
				}
				if ts.TalkSessionID != talkSessionID.UUID() {
					return errors.New("トークセッションIDが一致しません")
				}
				location := talksession.NewLocation(
					talkSessionID,
					ts.Latitude.(float64),
					ts.Longitude.(float64),
					ts.City.String,
					ts.Prefecture.String,
				)
				assert.Equal(t, location.ToGeographyText(), "POINT(30.000000 30.000000)")

				return nil
			},
			WantErr: false,
		},
	}

	txtest.RunTransactionalTests(t, dbManager, initData, testCases)
}
