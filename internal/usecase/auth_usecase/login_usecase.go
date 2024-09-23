package auth_usecase

import (
	"context"
	"net/http"
	"net/url"

	"github.com/neko-dream/server/internal/domain/model/auth"
	cookie_utils "github.com/neko-dream/server/pkg/cookie"
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
		Cookie      string
	}

	authLoginInteractor struct {
		auth.AuthService
	}
)

func (a *authLoginInteractor) Execute(ctx context.Context, input AuthLoginInput) (AuthLoginOutput, error) {
	authUrl, state, err := a.GetAuthURL(ctx, input.Provider)
	if err != nil {
		return AuthLoginOutput{}, err
	}

	stateCookie := http.Cookie{
		Name:     "state",
		Value:    state,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost",
	}
	redirectURLCookie := http.Cookie{
		Name:     "redirect_url",
		Value:    input.RedirectURL,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost",
	}
	encoded := cookie_utils.EncodeCookies([]*http.Cookie{
		&stateCookie,
		&redirectURLCookie,
	})

	return AuthLoginOutput{
		RedirectURL: authUrl,
		Cookie:      encoded,
	}, nil
}

func NewAuthLoginUseCase() AuthLoginUseCase {
	return &authLoginInteractor{}
}
