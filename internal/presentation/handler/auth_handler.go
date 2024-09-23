package handler

import (
	"context"
	"errors"
	"net/url"

	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/auth_usecase"
)

type authHandler struct {
	auth_usecase.AuthCallbackUseCase
}

// AuthLogin implements oas.AuthHandler.
func (a *authHandler) AuthLogin(ctx context.Context, params oas.AuthLoginParams) (*oas.AuthLoginFound, error) {
	panic("unimplemented")
}

// OAuthCallback implements oas.AuthHandler.
func (a *authHandler) OAuthCallback(ctx context.Context, params oas.OAuthCallbackParams) (*oas.OAuthCallbackFound, error) {
	if params.CookieState.Value != params.QueryState.Value {
		res := new(oas.OAuthCallbackFound)
		return res, errors.New("invalid state")
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

func NewAuthHandler() oas.AuthHandler {
	return &authHandler{}
}
