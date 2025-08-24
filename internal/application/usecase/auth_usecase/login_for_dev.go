package auth_usecase

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	organizationService "github.com/neko-dream/server/internal/domain/service/organization"
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
		// 組織コード（組織ログインの場合）
		OrganizationCode *string
	}

	LoginForDevOutput struct {
		Token string
	}

	loginForDevInteractor struct {
		*db.DBManager
		*config.Config
		service.AuthenticationService
		session.SessionRepository
		session.SessionService
		session.TokenManager
		organizationService organizationService.OrganizationService
	}
)

func NewLoginForDev(
	tm *db.DBManager,
	config *config.Config,
	authService service.AuthenticationService,
	sessionRepository session.SessionRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
	organizationService organizationService.OrganizationService,
) LoginForDev {
	return &loginForDevInteractor{
		DBManager:             tm,
		Config:                config,
		AuthenticationService: authService,
		SessionRepository:     sessionRepository,
		SessionService:        sessionService,
		TokenManager:          tokenManager,
		organizationService:   organizationService,
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
		newUser, err := a.AuthenticationService.Authenticate(ctx, "dev", input.Subject)
		if err != nil {
			return err
		}
		if newUser != nil {
			if err := a.SessionService.DeactivateUserSessions(ctx, newUser.UserID()); err != nil {
				utils.HandleError(ctx, err, "failed to deactivate user sessions")
			}
		}

		// 組織コードから組織IDを解決
		organizationID, err := a.organizationService.ResolveOrganizationIDFromCode(ctx, input.OrganizationCode)
		if err != nil {
			utils.HandleError(ctx, err, "ResolveOrganizationIDFromCode")
			return errtrace.Wrap(err)
		}

		// セッション作成
		var sess *session.Session
		if organizationID != nil {
			sess = session.NewSessionWithOrganization(
				shared.NewUUID[session.Session](),
				newUser.UserID(),
				newUser.Provider(),
				session.SESSION_ACTIVE,
				*session.NewExpiresAt(ctx),
				clock.Now(ctx),
				organizationID,
			)
		} else {
			sess = session.NewSession(
				shared.NewUUID[session.Session](),
				newUser.UserID(),
				newUser.Provider(),
				session.SESSION_ACTIVE,
				*session.NewExpiresAt(ctx),
				clock.Now(ctx),
			)
		}

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
