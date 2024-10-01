package user

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

type (
	RegisterUserUseCase interface {
		Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error)
	}

	RegisterUserInput struct {
		UserID        shared.UUID[user.User]
		DisplayID     string  // ユーザーの表示用ID
		DisplayName   string  // ユーザーの表示名
		PictureURL    *string // ユーザーのプロフィール画像
		YearOfBirth   *int    // ユーザーの生年
		Gender        *string // ユーザーの性別
		Municipality  *string // ユーザーの住んでいる市町村
		Occupation    *string // ユーザーの職業
		HouseholdSize *int    // ユーザーの世帯人数
	}

	RegisterUserOutput struct {
		DisplayID   string // ユーザーの表示用ID
		DisplayName string // ユーザーの表示名
	}

	registerUserInteractor struct {
		*db.DBManager
		userRep     user.UserRepository
		userService user.UserService
	}
)

func NewRegisterUserUseCase(
	tm *db.DBManager,
	userRep user.UserRepository,
	userService user.UserService,
) RegisterUserUseCase {
	return &registerUserInteractor{
		DBManager:   tm,
		userRep:     userRep,
		userService: userService,
	}
}

func (i *registerUserInteractor) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	err := i.ExecTx(ctx, func(ctx context.Context) error {
		// ユーザーの存在を確認
		foundUser, err := i.userRep.FindByID(ctx, input.UserID)
		if err != nil {
			utils.HandleError(ctx, err, "UserRepository.FindByID")
			return messages.UserNotFoundError
		}
		if foundUser == nil {
			return messages.UserNotFoundError
		}

		// ユーザーの表示用IDが重複していないかチェック
		duplicated, err := i.userService.DisplayIDCheckDuplicate(ctx, input.DisplayID)
		if err != nil {
			utils.HandleError(ctx, err, "UserService.DisplayIDCheckDuplicate")
			return messages.UserDisplayIDAlreadyExistsError
		}
		if duplicated {
			return messages.UserDisplayIDAlreadyExistsError
		}

		// ユーザー名と表示名を設定
		if err := foundUser.SetDisplayID(input.DisplayID); err != nil {
			utils.HandleError(ctx, err, "User.SetDisplayID")
			return messages.UserDisplayIDAlreadyExistsError
		}
		foundUser.ChangeName(input.DisplayName)

		// デモグラ情報を設定
		foundUser.SetDemographics(user.NewUserDemographics(
			user.NewYearOfBirth(input.YearOfBirth),
			user.NewOccupation(input.Occupation),
			lo.ToPtr(user.NewGender(input.Gender)),
			user.NewMunicipality(input.Municipality),
			user.NewHouseholdSize(input.HouseholdSize),
		))

		if err := i.userRep.Update(ctx, *foundUser); err != nil {
			utils.HandleError(ctx, err, "UserRepository.Update")
			return messages.UserUpdateError
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &RegisterUserOutput{
		DisplayID:   input.DisplayID,
		DisplayName: input.DisplayName,
	}, nil
}
