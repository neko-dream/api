package user_command

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	Edit interface {
		Execute(context.Context, EditInput) (*EditOutput, error)
	}

	EditInput struct {
		UserID        shared.UUID[user.User]
		DisplayName   *string               // ユーザーの表示名
		Icon          *multipart.FileHeader // ユーザーのアイコン
		DeleteIcon    bool                  // アイコンを削除するかどうか
		YearOfBirth   *int                  // ユーザーの生年
		Gender        *string               // ユーザーの性別
		City          *string               // ユーザーの住んでいる市町村
		Occupation    *string               // ユーザーの職業
		HouseholdSize *int                  // ユーザーの世帯人数
		Prefecture    *string               // ユーザーの居住地の都道府県
	}

	EditOutput struct {
		DisplayID   string // ユーザーの表示用ID
		DisplayName string // ユーザーの表示名
		Cookie      *http.Cookie
	}

	EditHandler struct {
		*db.DBManager
		session.TokenManager
		conf        *config.Config
		userRep     user.UserRepository
		userService user.UserService
		sessService session.SessionService
	}
)

func NewEditHandler(
	dm *db.DBManager,
	tm session.TokenManager,
	conf *config.Config,
	userRep user.UserRepository,
	userService user.UserService,
	sessService session.SessionService,
) Edit {
	return &EditHandler{
		DBManager:    dm,
		TokenManager: tm,
		conf:         conf,
		userRep:      userRep,
		userService:  userService,
		sessService:  sessService,
	}
}

func (e *EditHandler) Execute(ctx context.Context, input EditInput) (*EditOutput, error) {
	ctx, span := otel.Tracer("user_command").Start(ctx, "EditHandler.Execute")
	defer span.End()

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
		foundUser.ChangeName(ctx, input.DisplayName)

		// アイコンがある場合は設定
		if err := foundUser.SetIconFile(ctx, input.Icon); err != nil {
			utils.HandleError(ctx, err, "User.SetIconFile")
			return messages.UserUpdateError
		}
		if input.DeleteIcon {
			foundUser.DeleteIcon()
		}

		if input.YearOfBirth != nil ||
			input.Gender != nil ||
			input.City != nil ||
			input.Occupation != nil ||
			input.HouseholdSize != nil ||
			input.Prefecture != nil {
			var demograID shared.UUID[user.UserDemographics]
			if foundUser.Demographics() != nil {
				demograID = foundUser.Demographics().UserDemographicsID()
			} else {
				demograID = shared.NewUUID[user.UserDemographics]()
			}

			// デモグラ情報を設定
			foundUser.SetDemographics(user.NewUserDemographics(
				ctx,
				demograID,
				input.YearOfBirth,
				input.Occupation,
				input.Gender,
				input.City,
				input.HouseholdSize,
				input.Prefecture,
			))
		}

		if err := e.userRep.Update(ctx, *foundUser); err != nil {
			utils.HandleError(ctx, err, "UserRepository.Update")
			return messages.UserUpdateError
		}
		// 再度ユーザー情報を取得
		foundUser, err = e.userRep.FindByID(ctx, input.UserID)
		if err != nil {
			utils.HandleError(ctx, err, "UserRepository.FindByID")
			return messages.UserNotFoundError
		}
		u = foundUser

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

	return &EditOutput{
		DisplayID:   *u.DisplayID(),
		DisplayName: *u.DisplayName(),
		Cookie:      &c,
	}, nil
}
