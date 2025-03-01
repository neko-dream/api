package repository_test

import (
	"errors"
	"testing"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/crypto"
	ci "github.com/neko-dream/server/internal/infrastructure/crypto"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/server/internal/test/txtest"
	"github.com/samber/lo"
)

func TestUserRepository_Create(t *testing.T) {
	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)

	// 暗号化の設定
	encryptor, err := ci.NewEncryptor(lo.ToPtr(config.Config{
		ENCRYPTION_VERSION: crypto.Version1,
		ENCRYPTION_SECRET:  "12345678901234567890123456789012", // テスト用の32バイトキー
	}))
	if err != nil {
		t.Fatalf("暗号化の初期化に失敗: %v", err)
	}

	type TestData struct {
		UserRepo user.UserRepository
	}

	initData := TestData{}

	userID := shared.NewUUID[user.User]()
	testCases := []*txtest.TransactionalTestCase[TestData]{
		{
			Name: "ユーザー作成と暗号化された情報の検証",
			SetupFn: func(ctx *txtest.TestContext[TestData]) error {
				ctx.Data.UserRepo = repository.NewUserRepository(
					dbManager,
					repository.NewImageRepositoryMock(),
					encryptor,
				)
				return nil
			},
			TestFn: func(ctx *txtest.TestContext[TestData]) error {
				usr := user.NewUser(
					userID,
					lo.ToPtr("test"),
					lo.ToPtr("test"),
					"test",
					"GOOGLE",
					nil,
				)
				// 人口統計情報を設定
				demographics := user.NewUserDemographic(
					ctx.Context,
					shared.NewUUID[user.UserDemographic](),
					lo.ToPtr(1990),   // 生年
					lo.ToPtr("正社員"),  // 職業
					lo.ToPtr("男性"),   // 性別
					lo.ToPtr("世田谷区"), // 都市
					lo.ToPtr(2),      // 世帯人数
					lo.ToPtr("東京都"),  // 都道府県
				)
				usr.SetDemographics(demographics)

				// ユーザー作成
				err := ctx.Data.UserRepo.Create(ctx.Context, usr)
				if err != nil {
					return err
				}
				err = ctx.Data.UserRepo.Update(ctx.Context, usr)
				if err != nil {
					return err
				}

				// ユーザー情報の取得と検証
				foundUser, err := ctx.Data.UserRepo.FindByID(ctx.Context, userID)
				if err != nil {
					return err
				}

				if foundUser == nil {
					return errors.New("ユーザーが見つかりません")
				}

				// 基本情報の検証
				if foundUser.UserID() != userID {
					return errors.New("ユーザーIDが一致しません")
				}
				if foundUser.DisplayName() == nil {
					return errors.New("DisplayNameが見つかりません")
				}
				if *foundUser.DisplayName() != "test" {
					return errors.New("DisplayNameが一致しません")
				}
				if foundUser.DisplayID() == nil {
					return errors.New("DisplayIDが見つかりません")
				}
				if *foundUser.DisplayID() != "test" {
					return errors.New("DisplayIDが一致しません")
				}

				// 人口統計情報の検証
				if foundUser.Demographics() == nil {
					return errors.New("人口統計情報が見つかりません")
				}

				demo := *foundUser.Demographics()
				if *demo.YearOfBirth() != 1990 {
					return errors.New("生年が一致しません")
				}
				if demo.YearOfBirth() == nil {
					return errors.New("生年が見つかりません")
				}

				if demo.Occupation() == nil {
					return errors.New("職業が見つかりません")
				}
				if *demo.Occupation() != user.OccupationFullTimeEmployee {
					return errors.New("職業が一致しません")
				}

				if demo.City() == nil {
					return errors.New("市区町村が見つかりません")
				}
				if *demo.City() != "世田谷区" {
					return errors.New("市区町村が一致しません")
				}

				if demo.Prefecture() == nil {
					return errors.New("都道府県が見つかりません")
				}
				if *demo.Prefecture() != "東京都" {
					return errors.New("都道府県が一致しません")
				}

				if demo.HouseholdSize() == nil {
					return errors.New("世帯人数が見つかりません")
				}
				if *demo.HouseholdSize() != 2 {
					return errors.New("世帯人数が一致しません")
				}

				if demo.Gender() == nil {
					return errors.New("性別が見つかりません")
				}
				if *demo.Gender() != user.GenderMale {
					return errors.New("性別が一致しません")
				}

				return nil
			},
			WantErr: false,
		},
	}

	txtest.RunTransactionalTests(t, dbManager, initData, testCases)
}
