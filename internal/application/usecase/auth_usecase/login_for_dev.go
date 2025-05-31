package auth_usecase

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
	LoginForDev interface {
		Execute(context.Context, LoginForDevInput) (LoginForDevOutput, error)
	}

	LoginForDevInput struct {
		Subject string
	}

	LoginForDevOutput struct {
		Token string
	}

	loginForDevInteractor struct {
		*db.DBManager
		*config.Config
		service.AuthService
		session.SessionRepository
		session.SessionService
		session.TokenManager
	}
)

func NewLoginForDev(
	tm *db.DBManager,
	config *config.Config,
	authService service.AuthService,
	sessionRepository session.SessionRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
) LoginForDev {
	return &loginForDevInteractor{
		DBManager:         tm,
		Config:            config,
		AuthService:       authService,
		SessionRepository: sessionRepository,
		SessionService:    sessionService,
		TokenManager:      tokenManager,
	}
}

func (a *loginForDevInteractor) Execute(ctx context.Context, input LoginForDevInput) (LoginForDevOutput, error) {
	ctx, span := otel.Tracer("auth_usecase").Start(ctx, "loginForDevInteractor.Execute")
	defer span.End()

	if a.Config.Env == config.PROD {
		utils.HandleError(ctx, errtrace.New("このエンドポイントは開発でのみ有効です。"), "failed to login for dev")
		return LoginForDevOutput{}, errtrace.New("このエンドポイントは開発でのみ有効です。")
	}

	var (
		tok string
	)

	if err := a.ExecTx(ctx, func(ctx context.Context) error {
		newUser, err := a.AuthService.Authenticate(ctx, "dev", input.Subject)
		if err != nil {
			return err
		}
		if newUser != nil {
			if err := a.SessionService.DeactivateUserSessions(ctx, newUser.UserID()); err != nil {
				utils.HandleError(ctx, err, "failed to deactivate user sessions")
			}
		}

		sess := session.NewSession(
			shared.NewUUID[session.Session](),
			newUser.UserID(),
			newUser.Provider(),
			session.SESSION_ACTIVE,
			*session.NewExpiresAt(ctx),
			clock.Now(ctx),
		)

		if _, err := a.SessionRepository.Create(ctx, *sess); err != nil {
			utils.HandleError(ctx, err, "failed to create session")
			return errtrace.Wrap(err)
		}

		token, err := a.TokenManager.Generate(ctx, *newUser, sess.SessionID())
		if err != nil {
			utils.HandleError(ctx, err, "failed to generate token")
			return errtrace.Wrap(err)
		}

		tok = token

		return nil
	}); err != nil {
		return LoginForDevOutput{}, errtrace.Wrap(err)
	}

	return LoginForDevOutput{
		Token: tok,
	}, nil
}
