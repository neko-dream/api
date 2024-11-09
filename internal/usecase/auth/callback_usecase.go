package auth_usecase

import (
	"context"
	"net/http"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/db"
	cookie_utils "github.com/neko-dream/server/pkg/cookie"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	AuthCallbackUseCase interface {
		Execute(ctx context.Context, input CallbackInput) (CallbackOutput, error)
	}

	CallbackInput struct {
		Provider string
		Code     string
		Remember bool
	}

	CallbackOutput struct {
		Cookie string
	}

	authCallbackInteractor struct {
		*db.DBManager
		*config.Config
		auth.AuthService
		session.SessionRepository
		session.SessionService
		session.TokenManager
	}
)

func NewAuthCallbackUseCase(
	tm *db.DBManager,
	config *config.Config,
	authService auth.AuthService,
	sessionRepository session.SessionRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
) AuthCallbackUseCase {
	return &authCallbackInteractor{
		DBManager:         tm,
		Config:            config,
		AuthService:       authService,
		SessionRepository: sessionRepository,
		SessionService:    sessionService,
		TokenManager:      tokenManager,
	}
}

func (u *authCallbackInteractor) Execute(ctx context.Context, input CallbackInput) (CallbackOutput, error) {
	var c http.Cookie
	if err := u.ExecTx(ctx, func(ctx context.Context) error {
		user, err := u.AuthService.Authenticate(ctx, input.Provider, input.Code)
		if err != nil {
			utils.HandleError(ctx, err, "failed to authenticate")
			return errtrace.Wrap(err)
		}
		if user != nil {
			if err := u.SessionService.DeactivateUserSessions(ctx, user.UserID()); err != nil {
				utils.HandleError(ctx, err, "failed to deactivate user sessions")
			}
		}

		sess := session.NewSession(
			shared.NewUUID[session.Session](),
			user.UserID(),
			user.Provider(),
			session.SESSION_ACTIVE,
			*session.NewExpiresAt(ctx),
			time.Now(),
		)

		if _, err := u.SessionRepository.Create(ctx, *sess); err != nil {
			utils.HandleError(ctx, err, "failed to create session")
			return errtrace.Wrap(err)
		}

		token, err := u.TokenManager.Generate(ctx, *user, sess.SessionID())
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
			Domain:   u.DOMAIN,
			MaxAge:   60 * 60 * 24 * 7,
		}

		c = cookie
		return nil
	}); err != nil {
		return CallbackOutput{}, err
	}

	return CallbackOutput{
		Cookie: cookie_utils.EncodeCookies([]*http.Cookie{&c}),
	}, nil
}
