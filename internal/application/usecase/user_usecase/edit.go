package user_usecase

import (
	"context"
	"mime/multipart"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/service"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/presentation/oas"
	"github.com/neko-dream/api/pkg/utils"
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
		Email         *string               // ユーザーのメールアドレス
		DeleteIcon    bool                  // アイコンを削除するかどうか
		DateOfBirth   *int                  // ユーザーの生年
		Gender        *string               // ユーザーの性別
		City          *string               // ユーザーの住んでいる市区町村
		Occupation    *string               // ユーザーの職業
		HouseholdSize *int                  // ユーザーの世帯人数
		Prefecture    *string               // ユーザーの居住地の都道府県
	}

	EditOutput struct {
		DisplayID   string  // ユーザーの表示用ID
		DisplayName string  // ユーザーの表示名
		Token       string  // ユーザーのトークン
		IconURL     *string // ユーザーのアイコンURL
	}

	editHandler struct {
		*db.DBManager
		session.TokenManager
		conf               *config.Config
		userRep            user.UserRepository
		userService        user.UserService
		sessService        session.SessionService
		profileIconService service.ProfileIconService
	}
)

func (o *EditOutput) ToResponse() oas.User {
	return oas.User{
		DisplayID:   o.DisplayID,
		DisplayName: o.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](o.IconURL),
	}
}

func NewEditHandler(
	dm *db.DBManager,
	tm session.TokenManager,
	conf *config.Config,
	userRep user.UserRepository,
	userService user.UserService,
	sessService session.SessionService,
	profileIconService service.ProfileIconService,
) Edit {
	return &editHandler{
		DBManager:          dm,
		TokenManager:       tm,
		conf:               conf,
		userRep:            userRep,
		userService:        userService,
		sessService:        sessService,
		profileIconService: profileIconService,
	}
}

func (e *editHandler) Execute(ctx context.Context, input EditInput) (*EditOutput, error) {
	ctx, span := otel.Tracer("user_command").Start(ctx, "EditHandler.Execute")
	defer span.End()

	var (
		u     *user.User
		token string
	)

	err := e.ExecTx(ctx, func(ctx context.Context) error {
		// ユーザーの存在を確認
		foundUser, err := e.userRep.FindByID(ctx, input.UserID)
		if err != nil || foundUser == nil {
			utils.HandleError(ctx, err, "UserRepository.FindByID")
			return messages.UserNotFoundError
		}
		// ユーザーの表示名を変更
		foundUser.ChangeName(ctx, input.DisplayName)

		// アイコンがある場合はアップロード
		if input.Icon != nil {
			url, err := e.profileIconService.UploadProfileIcon(ctx, foundUser.UserID(), input.Icon)
			if err != nil {
				utils.HandleError(ctx, err, "SetIcon")
				return err
			}
			foundUser.ChangeIconURL(url)
		}
		if input.DeleteIcon {
			foundUser.DeleteIcon()
		}

		if input.Email != nil {
			// メールアドレスを変更
			foundUser.ChangeEmail(*input.Email)
			foundUser.SetEmailVerified(false)
		}
		var demograID shared.UUID[user.UserDemographic]
		if foundUser.Demographics() != nil {
			demograID = foundUser.Demographics().ID()
		} else {
			demograID = shared.NewUUID[user.UserDemographic]()
		}

		// デモグラ情報を設定
		foundUser.SetDemographics(user.NewUserDemographic(
			ctx,
			demograID,
			input.DateOfBirth,
			input.Gender,
			input.City,
			input.Prefecture,
		))

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

		tokenTmp, err := e.TokenManager.Generate(ctx, *foundUser, sess.SessionID())
		if err != nil {
			utils.HandleError(ctx, err, "failed to generate token")
			return messages.TokenGenerateError
		}

		token = tokenTmp
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &EditOutput{
		DisplayID:   *u.DisplayID(),
		DisplayName: *u.DisplayName(),
		Token:       token,
		IconURL:     u.IconURL(),
	}, nil
}
