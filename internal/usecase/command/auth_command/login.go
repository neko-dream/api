package auth_command

import (
	"context"
	"net/http"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	AuthLogin interface {
		Execute(context.Context, AuthLoginInput) (AuthLoginOutput, error)
	}

	AuthLoginInput struct {
		RedirectURL string
		Provider    string
	}

	AuthLoginOutput struct {
		RedirectURL string
		Cookies     []*http.Cookie
	}

	authLoginInteractor struct {
		*db.DBManager
		*config.Config
		service.AuthService
		authProviderFactory auth.AuthProviderFactory
	}
)

func NewAuthLogin(
	tm *db.DBManager,
	config *config.Config,
	authService service.AuthService,
	authProviderFactory auth.AuthProviderFactory,
) AuthLogin {
	return &authLoginInteractor{
		DBManager:           tm,
		Config:              config,
		AuthService:         authService,
		authProviderFactory: authProviderFactory,
	}
}

func (a *authLoginInteractor) Execute(ctx context.Context, input AuthLoginInput) (AuthLoginOutput, error) {
	ctx, span := otel.Tracer("auth_usecase").Start(ctx, "authLoginInteractor.Execute")
	defer span.End()

	var (
		s  string
		au string
	)

	if err := a.ExecTx(ctx, func(ctx context.Context) error {
		provider, err := a.authProviderFactory.NewAuthProvider(ctx, input.Provider)
		if err != nil {
			utils.HandleError(ctx, err, "NewAuthProvider")
			return errtrace.Wrap(err)
		}

		state, err := a.GenerateState(ctx)
		if err != nil {
			utils.HandleError(ctx, err, "GenerateState")
			return errtrace.Wrap(err)
		}

		s = state
		au = provider.GetAuthorizationURL(ctx, state)
		return nil
	}); err != nil {
		return AuthLoginOutput{}, errtrace.Wrap(err)
	}
	stateCookie := http.Cookie{
		Name:     "state",
		Value:    s,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   a.DOMAIN,
	}
	redirectURLCookie := http.Cookie{
		Name:     "redirect_url",
		Value:    input.RedirectURL,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   a.DOMAIN,
	}

	return AuthLoginOutput{
		RedirectURL: au,
		Cookies:     []*http.Cookie{&stateCookie, &redirectURLCookie},
	}, nil
}
