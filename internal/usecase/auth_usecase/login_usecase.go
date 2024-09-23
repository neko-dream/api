package auth_usecase

import (
	"context"
	"net/http"
	"net/url"

	"github.com/neko-dream/server/internal/domain/model/auth"
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
		auth.AuthService
	}
)

func (a *authLoginInteractor) Execute(ctx context.Context, input AuthLoginInput) (AuthLoginOutput, error) {
	authUrl, state, err := a.AuthService.GetAuthURL(ctx, input.Provider)
	if err != nil {
		return AuthLoginOutput{}, err
	}

	stateCookie := http.Cookie{
		Name:     "state",
		Value:    state,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost",
	}
	redirectURLCookie := http.Cookie{
		Name:     "redirect_url",
		Value:    input.RedirectURL,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost",
	}

	return AuthLoginOutput{
		RedirectURL: authUrl,
		Cookies:     []*http.Cookie{&stateCookie, &redirectURLCookie},
	}, nil
}

func NewAuthLoginUseCase(
	authService auth.AuthService,
) AuthLoginUseCase {
	return &authLoginInteractor{
		AuthService: authService,
	}
}
