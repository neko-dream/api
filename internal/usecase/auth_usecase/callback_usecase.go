package auth_usecase

import (
	"context"
	"net/http"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
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

	authCallbackUseCase struct {
		auth.AuthService
		session.SessionRepository
		session.SessionService
		session.TokenManager
	}
)

func NewAuthCallbackUseCase(
	authService auth.AuthService,
	sessionRepository session.SessionRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
) AuthCallbackUseCase {
	return &authCallbackUseCase{
		AuthService:       authService,
		SessionRepository: sessionRepository,
		SessionService:    sessionService,
		TokenManager:      tokenManager,
	}
}

func (u *authCallbackUseCase) Execute(ctx context.Context, input CallbackInput) (CallbackOutput, error) {
	user, err := u.AuthService.Authenticate(ctx, input.Provider, input.Code)
	if err != nil {
		return CallbackOutput{}, err
	}

	if err := u.SessionService.DeactivateUserSessions(ctx, user.UserID()); err != nil {
		return CallbackOutput{}, err
	}

	sess := session.NewSession(
		shared.NewUUID[session.Session](),
		user.UserID(),
		user.Provider(),
		session.SESSION_ACTIVE,
		*session.NewExpiresAt(),
	)

	if _, err := u.SessionRepository.Create(ctx, *sess); err != nil {
		return CallbackOutput{}, err
	}

	token, err := u.TokenManager.GenerateToken(ctx, user.UserID(), sess.SessionID())
	if err != nil {
		return CallbackOutput{}, err
	}

	cookie := http.Cookie{
		Name:     "SessionId",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost",
	}
	if input.Remember {
		cookie.MaxAge = 60 * 60 * 24 * 30
	}
	encoded := cookie_utils.EncodeCookies([]*http.Cookie{&cookie})
	return CallbackOutput{
		Cookie: encoded,
	}, nil
}
