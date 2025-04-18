package user_command

import (
	"context"
	"mime/multipart"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/service"
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
		SessionID   shared.UUID[session.Session] // セッションID
		UserID      shared.UUID[user.User]
		DisplayID   string                // ユーザーの表示用ID
		DisplayName string                // ユーザーの表示名
		Icon        *multipart.FileHeader // ユーザーのアイコン
		DateOfBirth *int                  // ユーザーの生年
		Gender      *string               // ユーザーの性別
		City        *string               // ユーザーの住んでいる市町村
		Prefecture  *string               // ユーザーの住んでいる都道府県
	}

	RegisterOutput struct {
		DisplayID   string  // ユーザーの表示用ID
		DisplayName string  // ユーザーの表示名
		IconURL     *string // ユーザーのアイコンURL
		Token       string  // ユーザーのトークン
	}

	registerHandler struct {
		*db.DBManager
		session.TokenManager
		conf               *config.Config
		userRep            user.UserRepository
		userService        user.UserService
		sessService        session.SessionService
		profileIconService service.ProfileIconService
	}
)

func NewRegisterHandler(
	dm *db.DBManager,
	tm session.TokenManager,
	conf *config.Config,
	userRep user.UserRepository,
	userService user.UserService,
	sessService session.SessionService,
	profileIconService service.ProfileIconService,
) Register {
	return &registerHandler{
		DBManager:          dm,
		TokenManager:       tm,
		conf:               conf,
		userRep:            userRep,
		userService:        userService,
		sessService:        sessService,
		profileIconService: profileIconService,
	}
}

func (i *registerHandler) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	ctx, span := otel.Tracer("user_command").Start(ctx, "registerHandler.Execute")
	defer span.End()

	var (
		iconURL *string
		token   string
	)

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

		if input.Icon != nil {
			url, err := i.profileIconService.UploadProfileIcon(ctx, foundUser.UserID(), input.Icon)
			if err != nil {
				return err
			}
			foundUser.ChangeIconURL(url)
		}

		// デモグラ情報を設定
		foundUser.SetDemographics(user.NewUserDemographic(
			ctx,
			shared.NewUUID[user.UserDemographic](),
			input.DateOfBirth,
			input.Gender,
			input.City,
			input.Prefecture,
		))

		if err := i.userRep.Update(ctx, *foundUser); err != nil {
			utils.HandleError(ctx, err, "UserRepository.Update")
			return messages.UserUpdateError
		}

		if foundUser.IconURL() != nil {
			iconURL = foundUser.IconURL()
		}

		tokenTmp, err := i.TokenManager.Generate(ctx, *foundUser, input.SessionID)
		if err != nil {
			utils.HandleError(ctx, err, "failed to generate token")
			return errtrace.Wrap(err)
		}
		token = tokenTmp
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{
		DisplayID:   input.DisplayID,
		DisplayName: input.DisplayName,
		IconURL:     iconURL,
		Token:       token,
	}, nil
}
