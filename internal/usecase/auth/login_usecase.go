package auth_usecase

import (
	"context"
	"net/http"
	"net/url"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/db"
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
		auth.AuthService
	}
)

func (a *authLoginInteractor) Execute(ctx context.Context, input AuthLoginInput) (AuthLoginOutput, error) {
	var (
		s  string
		au *url.URL
	)

	if err := a.ExecTx(ctx, func(ctx context.Context) error {
		authURL, state, err := a.AuthService.GetAuthURL(ctx, input.Provider)
		if err != nil {
			return errtrace.Wrap(err)
		}

		s = state
		au = authURL
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

func NewAuthLoginUseCase(
	tm *db.DBManager,
	config *config.Config,
	authService auth.AuthService,
) AuthLoginUseCase {
	return &authLoginInteractor{
		DBManager:   tm,
		Config:      config,
		AuthService: authService,
	}
}
