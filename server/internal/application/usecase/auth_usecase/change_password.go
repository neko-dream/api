package auth_usecase

import (
	"context"

	"github.com/neko-dream/api/internal/domain/messages"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type ChangePassword interface {
	Execute(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, error)
}

type ChangePasswordInput struct {
	UserID shared.UUID[user.User]
	// SessionID 変更操作を行った現在のセッション。これ以外のセッションは無効にする。
	SessionID   shared.UUID[session.Session]
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
	txRunner loginTransactionRunner
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
		txRunner:            dbManager,
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

	// 新しいパスワードを登録し、漏洩したパスワードで作られた可能性のある他のセッションを無効にする
	if err := c.txRunner.ExecTx(ctx, func(ctx context.Context) error {
		if err := c.passwordAuthManager.UpdatePassword(ctx, input.UserID, input.NewPassword); err != nil {
			return err
		}
		return c.deactivateOtherSessions(ctx, input.UserID, input.SessionID)
	}); err != nil {
		return nil, err
	}

	return &ChangePasswordOutput{Success: true}, nil
}

func (c *changePasswordInteractor) deactivateOtherSessions(ctx context.Context, userID shared.UUID[user.User], currentSessionID shared.UUID[session.Session]) error {
	sessions, err := c.SessionRepository.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	for _, sess := range sessions {
		if sess.SessionID() == currentSessionID || !sess.IsActive(ctx) {
			continue
		}
		sess.Deactivate(ctx)
		if _, err := c.SessionRepository.Update(ctx, sess); err != nil {
			return err
		}
	}
	return nil
}
