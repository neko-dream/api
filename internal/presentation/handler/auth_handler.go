package handler

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/command/auth_command"
	cookie_utils "github.com/neko-dream/server/pkg/cookie"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type authHandler struct {
	auth_command.AuthCallback
	auth_command.AuthLogin
	auth_command.Revoke
}

func NewAuthHandler(
	authLogin auth_command.AuthLogin,
	authCallback auth_command.AuthCallback,
	revoke auth_command.Revoke,
) oas.AuthHandler {
	return &authHandler{
		AuthLogin:    authLogin,
		AuthCallback: authCallback,
		Revoke:       revoke,
	}
}

// Authorize implements oas.AuthHandler.
func (a *authHandler) Authorize(ctx context.Context, params oas.AuthorizeParams) (oas.AuthorizeRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.Authorize")
	defer span.End()

	out, err := a.AuthLogin.Execute(ctx, auth_command.AuthLoginInput{
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
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.OAuthCallback")
	defer span.End()

	// CookieStateとQueryStateが一致しているか確認
	if params.CookieState != params.QueryState {
		return nil, messages.InvalidStateError
	}

	output, err := a.AuthCallback.Execute(ctx, auth_command.CallbackInput{
		Provider: params.Provider,
		Code:     params.Code,
	})
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
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.OAuthTokenRevoke")
	defer span.End()

	claim := session.GetSession(ctx)
	sessID, err := claim.SessionID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	out, err := a.Revoke.Execute(ctx, auth_command.RevokeInput{
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
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.OAuthTokenInfo")
	defer span.End()

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
