package repository_test

import (
	"errors"
	"testing"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/server/internal/test/txtest"
	"github.com/samber/lo"
)

func TestUserRepository_Create(t *testing.T) {
	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)

	type TestData struct {
		UserRepo user.UserRepository
	}

	initData := TestData{
		UserRepo: repository.NewUserRepository(dbManager, repository.NewImageRepositoryMock()),
	}
	userID := shared.NewUUID[user.User]()
	testCases := []*txtest.TransactionalTestCase[TestData]{
		{
			Name: "ユーザー作成ができる",
			TestFn: func(ctx *txtest.TestContext[TestData]) error {
				user := user.NewUser(
					userID,
					lo.ToPtr("test"),
					lo.ToPtr("test"),
					"test",
					"GOOGLE",
					nil,
				)
				err := ctx.Data.UserRepo.Create(ctx, user)
				if err != nil {
					return err
				}

				usr, err := dbManager.GetQueries(ctx).GetUserByID(ctx, userID.UUID())
				if err != nil {
					return err
				}
				if usr.UserID != userID.UUID() {
					return errors.New("ユーザーIDが一致しません")
				}

				err = ctx.Data.UserRepo.Update(ctx, user)
				if err != nil {
					return err
				}

				usr, err = dbManager.GetQueries(ctx).GetUserByID(ctx, userID.UUID())
				if err != nil {
					return err
				}
				if usr.UserID != userID.UUID() {
					return errors.New("ユーザーIDが一致しません")
				}

				return nil
			},
			WantErr: false,
		},
	}

	txtest.RunTransactionalTests(t, dbManager, initData, testCases)
}
