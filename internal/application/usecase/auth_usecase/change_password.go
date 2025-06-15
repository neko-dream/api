package auth_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	password_auth "github.com/neko-dream/server/internal/domain/model/auth/password"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type ChangePassword interface {
	Execute(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, error)
}

type ChangePasswordInput struct {
	UserID      shared.UUID[user.User]
	OldPassword string
	NewPassword string
}
type ChangePasswordOutput struct {
	Success bool
}

type changePasswordInteractor struct {
	passwordAuthManager password_auth.PasswordAuthManager
	userRepository      user.UserRepository
	cfg                 *config.Config
	session.TokenManager
	session.SessionRepository
	*db.DBManager
}

func NewChangePassword(
	passwordAuthManager password_auth.PasswordAuthManager,
	userRepository user.UserRepository,
	cfg *config.Config,
	tokenManager session.TokenManager,
	sessionRepository session.SessionRepository,
	dbManager *db.DBManager,
) ChangePassword {
	return &changePasswordInteractor{
		passwordAuthManager: passwordAuthManager,
		userRepository:      userRepository,
		cfg:                 cfg,
		TokenManager:        tokenManager,
		SessionRepository:   sessionRepository,
		DBManager:           dbManager,
	}
}

func (c *changePasswordInteractor) Execute(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, error) {
	ctx, span := otel.Tracer("auth_command").Start(ctx, "changePasswordInteractor.Execute")
	defer span.End()

	// ユーザーのパスワード認証情報を取得
	auth, err := c.passwordAuthManager.VerifyPassword(ctx, input.UserID, input.OldPassword)
	if err != nil {
		return nil, err
	}
	if !auth {
		return nil, messages.InvalidPasswordError
	}

	// 新しいパスワードを登録
	err = c.passwordAuthManager.UpdatePassword(ctx, input.UserID, input.NewPassword)
	if err != nil {
		return nil, err
	}

	return &ChangePasswordOutput{Success: true}, nil
}
