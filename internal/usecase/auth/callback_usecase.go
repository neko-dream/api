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
	authService auth.AuthService,
	sessionRepository session.SessionRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
) AuthCallbackUseCase {
	return &authCallbackInteractor{
		DBManager:         tm,
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
			return errtrace.Wrap(err)
		}
		if user != nil {
			if err := u.SessionService.DeactivateUserSessions(ctx, user.UserID()); err != nil {
				return errtrace.Wrap(err)
			}
		}

		sess := session.NewSession(
			shared.NewUUID[session.Session](),
			user.UserID(),
			user.Provider(),
			session.SESSION_ACTIVE,
			*session.NewExpiresAt(),
			time.Now(),
		)

		if _, err := u.SessionRepository.Create(ctx, *sess); err != nil {
			return errtrace.Wrap(err)
		}

		token, err := u.TokenManager.Generate(ctx, *user, sess.SessionID())
		if err != nil {
			return errtrace.Wrap(err)
		}
		cookie := http.Cookie{
			Name:     "SessionId",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Domain:   u.DOMAIN,
		}
		if input.Remember {
			cookie.MaxAge = 60 * 60 * 24 * 30
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
