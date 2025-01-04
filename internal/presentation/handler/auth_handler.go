package handler

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/presentation/oas"
	auth_usecase "github.com/neko-dream/server/internal/usecase/auth"
	cookie_utils "github.com/neko-dream/server/pkg/cookie"
	"github.com/neko-dream/server/pkg/utils"
)

type authHandler struct {
	auth_usecase.AuthCallbackUseCase
	auth_usecase.AuthLoginUseCase
	auth_usecase.RevokeUseCase
}

func NewAuthHandler(
	authCallbackUseCase auth_usecase.AuthCallbackUseCase,
	authLoginUseCase auth_usecase.AuthLoginUseCase,
	authRevokeUseCase auth_usecase.RevokeUseCase,
) oas.AuthHandler {
	return &authHandler{
		AuthCallbackUseCase: authCallbackUseCase,
		AuthLoginUseCase:    authLoginUseCase,
		RevokeUseCase:       authRevokeUseCase,
	}
}

// Authorize implements oas.AuthHandler.
func (a *authHandler) Authorize(ctx context.Context, params oas.AuthorizeParams) (oas.AuthorizeRes, error) {
	out, err := a.AuthLoginUseCase.Execute(ctx, auth_usecase.AuthLoginInput{
		RedirectURL: params.RedirectURL,
		Provider:    params.Provider,
	})
	if err != nil {
		return nil, err
	}

	headers := new(oas.AuthorizeFoundHeaders)
	headers.SetLocation(out.RedirectURL)
	headers.SetSetCookie(cookie_utils.EncodeCookies(out.Cookies))
	return headers, nil
}

// OAuthCallback implements oas.AuthHandler.
func (a *authHandler) OAuthCallback(ctx context.Context, params oas.OAuthCallbackParams) (oas.OAuthCallbackRes, error) {
	// CookieStateとQueryStateが一致しているか確認
	if params.CookieState != params.QueryState {
		return nil, messages.InvalidStateError
	}

	input := auth_usecase.CallbackInput{
		Provider: params.Provider,
		Code:     params.Code,
	}
	output, err := a.AuthCallbackUseCase.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	headers := new(oas.OAuthCallbackFoundHeaders)
	headers.SetCookie = output.Cookie
	// LoginでRedirectURLを設定しているためエラーは発生しない
	headers.Location = params.RedirectURL
	return headers, nil
}

// OAuthRevoke implements oas.AuthHandler.
func (a *authHandler) OAuthTokenRevoke(ctx context.Context) (oas.OAuthTokenRevokeRes, error) {
	claim := session.GetSession(ctx)
	sessID, err := claim.SessionID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	out, err := a.RevokeUseCase.Execute(ctx, auth_usecase.RevokeInput{
		SessID: sessID,
	})
	if err != nil {
		return nil, err
	}

	headers := new(oas.OAuthTokenRevokeNoContentHeaders)
	headers.SetSetCookie(cookie_utils.EncodeCookies(out.Cookies))
	return headers, nil
}

// OAuthTokenInfo implements oas.AuthHandler.
func (a *authHandler) OAuthTokenInfo(ctx context.Context) (oas.OAuthTokenInfoRes, error) {
	claim := session.GetSession(ctx)
	sessID, err := claim.SessionID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if claim.IsExpired(ctx) {
		return nil, messages.TokenExpiredError
	}

	return &oas.OAuthTokenInfoOK{
		Aud:         claim.Audience(),
		Iat:         claim.IssueAt().Format(time.RFC3339),
		Exp:         claim.ExpiresAt().Format(time.RFC3339),
		Iss:         claim.Issuer(),
		Sub:         claim.Sub,
		Jti:         sessID.String(),
		IsVerify:    claim.IsVerify,
		DisplayId:   utils.ToOpt[oas.OptString](claim.DisplayID),
		DisplayName: utils.ToOpt[oas.OptString](claim.DisplayName),
		IconURL:     utils.ToOpt[oas.OptString](claim.IconURL),
	}, nil
}
