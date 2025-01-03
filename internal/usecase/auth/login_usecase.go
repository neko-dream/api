package auth_usecase

import (
	"context"
	"net/http"
	"net/url"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	AuthLoginUseCase interface {
		Execute(context.Context, AuthLoginInput) (AuthLoginOutput, error)
	}

	AuthLoginInput struct {
		RedirectURL string
		Provider    string
	}

	AuthLoginOutput struct {
		RedirectURL *url.URL
		Cookies     []*http.Cookie
	}

	authLoginInteractor struct {
		*db.DBManager
		*config.Config
		service.AuthService
		authProviderFactory auth.AuthProviderFactory
	}
)

func NewAuthLoginUseCase(
	tm *db.DBManager,
	config *config.Config,
	authService service.AuthService,
	authProviderFactory auth.AuthProviderFactory,
) AuthLoginUseCase {
	return &authLoginInteractor{
		DBManager:           tm,
		Config:              config,
		AuthService:         authService,
		authProviderFactory: authProviderFactory,
	}
}

func (a *authLoginInteractor) Execute(ctx context.Context, input AuthLoginInput) (AuthLoginOutput, error) {
	var (
		s  string
		au *url.URL
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

		url, err := url.Parse(provider.GetAuthorizationURL(ctx, state))
		if err != nil {
			utils.HandleError(ctx, err, "url.Parse")
			return errtrace.Wrap(err)
		}

		s = state
		au = url
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
