package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/infrastructure/http/cookie"
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
	auth_command.LoginForDev

	cookie.CookieManager
}

func NewAuthHandler(
	authLogin auth_command.AuthLogin,
	authCallback auth_command.AuthCallback,
	revoke auth_command.Revoke,
	devLogin auth_command.LoginForDev,
	cookieManger cookie.CookieManager,
) oas.AuthHandler {
	return &authHandler{
		AuthLogin:     authLogin,
		AuthCallback:  authCallback,
		Revoke:        revoke,
		LoginForDev:   devLogin,
		CookieManager: cookieManger,
	}
}

// Authorize implements oas.AuthHandler.
func (a *authHandler) Authorize(ctx context.Context, params oas.AuthorizeParams) (oas.AuthorizeRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.Authorize")
	defer span.End()

	provider, err := params.Provider.MarshalText()
	if err != nil {
		return nil, err
	}

	out, err := a.AuthLogin.Execute(ctx, auth_command.AuthLoginInput{
		Provider: string(provider),
	})
	if err != nil {
		return nil, err
	}

	headers := new(oas.AuthorizeFoundHeaders)
	headers.SetLocation(out.RedirectURL)
	headers.SetSetCookie(cookie_utils.EncodeCookies(a.CookieManager.CreateAuthCookies(out.State, params.RedirectURL)))
	return headers, nil
}

func (a *authHandler) DevAuthorize(ctx context.Context, params oas.DevAuthorizeParams) (oas.DevAuthorizeRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.DevAuthorize")
	defer span.End()

	output, err := a.LoginForDev.Execute(ctx, auth_command.LoginForDevInput{
		Subject: params.ID,
	})
	if err != nil {
		return nil, err
	}

	headers := new(oas.DevAuthorizeFoundHeaders)
	headers.SetSetCookie(cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateSessionCookie(output.Token)}))
	headers.SetLocation(oas.NewOptString(params.RedirectURL))
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
	headers.SetSetCookie(
		cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateSessionCookie(output.Token)})[0],
	)
	// LoginでRedirectURLを設定しているためエラーは発生しない
	headers.SetLocation(params.RedirectURL)
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
	_, err = a.Revoke.Execute(ctx, auth_command.RevokeInput{
		SessID: sessID,
	})
	if err != nil {
		return nil, err
	}

	headers := new(oas.OAuthTokenRevokeNoContentHeaders)
	headers.SetSetCookie(cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateRevokeCookie()}))
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
		Aud:             claim.Audience(),
		Iat:             claim.IssueAt().Format(time.RFC3339),
		Exp:             claim.ExpiresAt().Format(time.RFC3339),
		Iss:             claim.Issuer(),
		Sub:             claim.Sub,
		Jti:             sessID.String(),
		IsVerify:        claim.IsVerify,
		IsEmailVerified: claim.IsEmailVerified,
		DisplayId:       utils.ToOpt[oas.OptString](claim.DisplayID),
		DisplayName:     utils.ToOpt[oas.OptString](claim.DisplayName),
		IconURL:         utils.ToOpt[oas.OptString](claim.IconURL),
	}, nil
}
