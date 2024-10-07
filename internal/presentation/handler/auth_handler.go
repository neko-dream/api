package handler

import (
	"context"
	"net/http"
	"net/url"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/presentation/oas"
	auth_usecase "github.com/neko-dream/server/internal/usecase/auth"
	http_utils "github.com/neko-dream/server/pkg/http"
)

type authHandler struct {
	auth_usecase.AuthCallbackUseCase
	auth_usecase.AuthLoginUseCase
}

func NewAuthHandler(
	authCallbackUseCase auth_usecase.AuthCallbackUseCase,
	authLoginUseCase auth_usecase.AuthLoginUseCase,
) oas.AuthHandler {
	return &authHandler{
		AuthCallbackUseCase: authCallbackUseCase,
		AuthLoginUseCase:    authLoginUseCase,
	}
}

// Authorize implements oas.AuthHandler.
func (a *authHandler) Authorize(ctx context.Context, params oas.AuthorizeParams) (*oas.AuthorizeFound, error) {
	out, err := a.AuthLoginUseCase.Execute(ctx, auth_usecase.AuthLoginInput{
		RedirectURL: params.RedirectURL,
		Provider:    params.Provider,
	})
	if err != nil {
		return nil, err
	}

	res := new(oas.AuthorizeFound)
	res.SetLocation(oas.NewOptURI(*out.RedirectURL))
	w := http_utils.GetHTTPResponse(ctx)
	for _, c := range out.Cookies {
		http.SetCookie(w, c)
	}

	return res, nil
}

// OAuthCallback implements oas.AuthHandler.
func (a *authHandler) OAuthCallback(ctx context.Context, params oas.OAuthCallbackParams) (*oas.OAuthCallbackFound, error) {
	if params.CookieState.Value != params.QueryState.Value {
		res := new(oas.OAuthCallbackFound)
		return res, messages.InvalidStateError
	}

	input := auth_usecase.CallbackInput{
		Provider: params.Provider,
		Code:     params.Code.Value,
	}

	output, err := a.AuthCallbackUseCase.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	res := new(oas.OAuthCallbackFound)
	res.SetCookie = oas.NewOptString(output.Cookie)
	// LoginでRedirectURLを設定しているためエラーは発生しない
	loc, _ := url.Parse(params.RedirectURL)
	res.Location = oas.NewOptURI(*loc)

	return res, nil
}
