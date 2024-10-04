package user

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

type (
	EditUserUseCase interface {
		Execute(ctx context.Context, input EditUserInput) (*EditUserOutput, error)
	}

	EditUserInput struct {
		UserID        shared.UUID[user.User]
		DisplayName   *string               // ユーザーの表示名
		Icon          *multipart.FileHeader // ユーザーのアイコン
		YearOfBirth   *int                  // ユーザーの生年
		Gender        *string               // ユーザーの性別
		Municipality  *string               // ユーザーの住んでいる市町村
		Occupation    *string               // ユーザーの職業
		HouseholdSize *int                  // ユーザーの世帯人数
	}

	EditUserOutput struct {
		DisplayID   string // ユーザーの表示用ID
		DisplayName string // ユーザーの表示名
		Cookie      *http.Cookie
	}

	editUserInteractor struct {
		*db.DBManager
		session.TokenManager
		conf        *config.Config
		userRep     user.UserRepository
		userService user.UserService
		sessService session.SessionService
	}
)

func NewEditUserUseCase(
	dm *db.DBManager,
	tm session.TokenManager,
	conf *config.Config,
	userRep user.UserRepository,
	userService user.UserService,
	sessService session.SessionService,
) EditUserUseCase {
	return &editUserInteractor{
		DBManager:    dm,
		TokenManager: tm,
		conf:         conf,
		userRep:      userRep,
		userService:  userService,
		sessService:  sessService,
	}
}

func (e *editUserInteractor) Execute(ctx context.Context, input EditUserInput) (*EditUserOutput, error) {
	var c http.Cookie
	var u *user.User

	err := e.ExecTx(ctx, func(ctx context.Context) error {
		// ユーザーの存在を確認
		foundUser, err := e.userRep.FindByID(ctx, input.UserID)
		if err != nil {
			utils.HandleError(ctx, err, "UserRepository.FindByID")
			return messages.UserNotFoundError
		}
		if foundUser == nil {
			return messages.UserNotFoundError
		}

		// ユーザー名と表示名を設定
		if input.Icon != nil {
			if err := foundUser.SetIconFile(ctx, input.Icon); err != nil {
				utils.HandleError(ctx, err, "User.SetIconFile")
				return messages.UserUpdateError
			}
		}
		if input.YearOfBirth != nil ||
			input.Gender != nil ||
			input.Municipality != nil ||
			input.Occupation != nil ||
			input.HouseholdSize != nil {
			var demograID shared.UUID[user.UserDemographics]
			if foundUser.Demographics() != nil {
				demograID = foundUser.Demographics().UserDemographicsID()
			} else {
				demograID = shared.NewUUID[user.UserDemographics]()
			}

			// デモグラ情報を設定
			foundUser.SetDemographics(user.NewUserDemographics(
				demograID,
				user.NewYearOfBirth(input.YearOfBirth),
				user.NewOccupation(input.Occupation),
				lo.ToPtr(user.NewGender(input.Gender)),
				user.NewMunicipality(input.Municipality),
				user.NewHouseholdSize(input.HouseholdSize),
			))
		}

		if err := e.userRep.Update(ctx, *foundUser); err != nil {
			utils.HandleError(ctx, err, "UserRepository.Update")
			return messages.UserUpdateError
		}

		sess, err := e.sessService.RefreshSession(ctx, input.UserID)
		if err != nil {
			utils.HandleError(ctx, err, "SessionService.RefreshSession")
			return messages.UserUpdateError
		}

		token, err := e.TokenManager.Generate(ctx, *foundUser, sess.SessionID())
		if err != nil {
			utils.HandleError(ctx, err, "failed to generate token")
			return messages.TokenGenerateError
		}
		cookie := http.Cookie{
			Name:     "SessionId",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Domain:   e.conf.DOMAIN,
			MaxAge:   60 * 60 * 24 * 7,
		}

		c = cookie
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &EditUserOutput{
		DisplayID:   *u.DisplayID(),
		DisplayName: *u.DisplayName(),
		Cookie:      &c,
	}, nil
}