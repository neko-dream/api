package auth_command

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	AuthCallback interface {
		Execute(ctx context.Context, input CallbackInput) (CallbackOutput, error)
	}

	CallbackInput struct {
		Provider string
		Code     string
		Remember bool
	}

	CallbackOutput struct {
		Token string
	}

	authCallbackInteractor struct {
		*db.DBManager
		*config.Config
		service.AuthService
		session.SessionRepository
		session.SessionService
		session.TokenManager
	}
)

func NewAuthCallback(
	tm *db.DBManager,
	config *config.Config,
	authService service.AuthService,
	sessionRepository session.SessionRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
) AuthCallback {
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
	ctx, span := otel.Tracer("auth_usecase").Start(ctx, "authCallbackInteractor.Execute")
	defer span.End()

	var tokenRes string
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
			clock.Now(ctx),
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

		tokenRes = token

		return nil
	}); err != nil {
		return CallbackOutput{}, err
	}

	return CallbackOutput{
		Token: tokenRes,
	}, nil
}
