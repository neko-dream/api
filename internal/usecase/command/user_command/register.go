package user_command

import (
	"context"
	"mime/multipart"
	"net/http"

	"braces.dev/errtrace"
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
	Register interface {
		Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error)
	}

	RegisterInput struct {
		UserID        shared.UUID[user.User]
		DisplayID     string                // ユーザーの表示用ID
		DisplayName   string                // ユーザーの表示名
		Icon          *multipart.FileHeader // ユーザーのアイコン
		YearOfBirth   *int                  // ユーザーの生年
		Gender        *string               // ユーザーの性別
		City          *string               // ユーザーの住んでいる市町村
		Occupation    *string               // ユーザーの職業
		HouseholdSize *int                  // ユーザーの世帯人数
		Prefecture    *string               // ユーザーの住んでいる都道府県
	}

	RegisterOutput struct {
		DisplayID   string  // ユーザーの表示用ID
		DisplayName string  // ユーザーの表示名
		IconURL     *string // ユーザーのアイコンURL
		Cookie      *http.Cookie
	}

	registerHandler struct {
		*db.DBManager
		session.TokenManager
		conf        *config.Config
		userRep     user.UserRepository
		userService user.UserService
		sessService session.SessionService
	}
)

func NewRegisterHandler(
	dm *db.DBManager,
	tm session.TokenManager,
	conf *config.Config,
	userRep user.UserRepository,
	userService user.UserService,
	sessService session.SessionService,
) Register {
	return &registerHandler{
		DBManager:    dm,
		TokenManager: tm,
		conf:         conf,
		userRep:      userRep,
		userService:  userService,
		sessService:  sessService,
	}
}

func (i *registerHandler) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	ctx, span := otel.Tracer("user_command").Start(ctx, "registerHandler.Execute")
	defer span.End()

	var c http.Cookie
	var iconURL *string

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

		foundUser.ChangeName(ctx, &input.DisplayName)

		if err := foundUser.SetIconFile(ctx, input.Icon); err != nil {
			utils.HandleError(ctx, err, "User.SetIconFile")
			return messages.UserUpdateError
		}

		// アイコンが指定されていればアイコンを設定
		if input.Icon != nil {
			if err := foundUser.SetIconFile(ctx, input.Icon); err != nil {
				utils.HandleError(ctx, err, "User.SetIconFile")
				return messages.UserUpdateError
			}
		}

		// デモグラ情報を設定
		foundUser.SetDemographics(user.NewUserDemographic(
			ctx,
			shared.NewUUID[user.UserDemographic](),
			input.YearOfBirth,
			input.Occupation,
			input.Gender,
			input.City,
			input.HouseholdSize,
			input.Prefecture,
		))

		if err := i.userRep.Update(ctx, *foundUser); err != nil {
			utils.HandleError(ctx, err, "UserRepository.Update")
			return messages.UserUpdateError
		}

		if foundUser.ProfileIconURL() != nil {
			iconURL = foundUser.ProfileIconURL()
		}

		sess, err := i.sessService.RefreshSession(ctx, input.UserID)
		if err != nil {
			utils.HandleError(ctx, err, "SessionService.RefreshSession")
			return messages.UserUpdateError
		}

		token, err := i.TokenManager.Generate(ctx, *foundUser, sess.SessionID())
		if err != nil {
			utils.HandleError(ctx, err, "failed to generate token")
			return errtrace.Wrap(err)
		}
		cookie := http.Cookie{
			Name:     "SessionId",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Domain:   i.conf.DOMAIN,
			MaxAge:   60 * 60 * 24 * 7,
		}

		c = cookie
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{
		DisplayID:   input.DisplayID,
		DisplayName: input.DisplayName,
		IconURL:     iconURL,
		Cookie:      &c,
	}, nil
}
