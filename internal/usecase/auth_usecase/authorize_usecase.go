package auth_usecase

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/infrastructure/db"
)

type (
	AuthorizeUseCase interface {
		Execute(context.Context, AuthLoginInput) (AuthLoginOutput, error)
	}

	AuthorizeInput struct {
		RedirectURL string
		Provider    string
	}

	AuthorizeOutput struct {
		RedirectURL *url.URL
		Cookies     []*http.Cookie
	}

	authorizeInteractor struct {
		*db.DBManager
		auth.AuthService
	}
)

func (a *authorizeInteractor) Execute(ctx context.Context, input AuthLoginInput) (AuthLoginOutput, error) {
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
	domain := os.Getenv("DOMAIN")
	stateCookie := http.Cookie{
		Name:     "state",
		Value:    s,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
	}
	redirectURLCookie := http.Cookie{
		Name:     "redirect_url",
		Value:    input.RedirectURL,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
	}

	return AuthLoginOutput{
		RedirectURL: au,
		Cookies:     []*http.Cookie{&stateCookie, &redirectURLCookie},
	}, nil
}

func NewAuthorizeUseCase(
	tm *db.DBManager,
	authService auth.AuthService,
) AuthorizeUseCase {
	return &authorizeInteractor{
		DBManager:   tm,
		AuthService: authService,
	}
}
