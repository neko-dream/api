package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/neko-dream/server/internal/application/usecase/auth_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/http/cookie"
	"github.com/neko-dream/server/internal/presentation/oas"
	cookie_utils "github.com/neko-dream/server/pkg/cookie"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type authHandler struct {
	auth_usecase.AuthCallback
	auth_usecase.AuthLogin
	auth_usecase.Revoke
	auth_usecase.LoginForDev
	auth_usecase.DetachAccount

	passwordLogin    auth_usecase.PasswordLogin
	passwordRegister auth_usecase.PasswordRegister
	changePassword   auth_usecase.ChangePassword

	cookie.CookieManager
}

func NewAuthHandler(
	authLogin auth_usecase.AuthLogin,
	authCallback auth_usecase.AuthCallback,
	revoke auth_usecase.Revoke,
	devLogin auth_usecase.LoginForDev,
	detachAccount auth_usecase.DetachAccount,

	login auth_usecase.PasswordLogin,
	register auth_usecase.PasswordRegister,
	changePassword auth_usecase.ChangePassword,

	cookieManger cookie.CookieManager,
) oas.AuthHandler {
	return &authHandler{
		AuthLogin:        authLogin,
		AuthCallback:     authCallback,
		Revoke:           revoke,
		LoginForDev:      devLogin,
		DetachAccount:    detachAccount,
		CookieManager:    cookieManger,
		passwordLogin:    login,
		passwordRegister: register,
		changePassword:   changePassword,
	}
}

// Authorize 認証プロバイダーの認可URLとstateを生成
func (a *authHandler) Authorize(ctx context.Context, params oas.AuthorizeParams) (oas.AuthorizeRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.Authorize")
	defer span.End()

	provider, err := params.Provider.MarshalText()
	if err != nil {
		return nil, err
	}
	var registrationURL *string
	if params.RegistrationURL.Set {
		registrationURL = lo.ToPtr(params.RegistrationURL.Value)
	}

	var organizationCode *string
	if params.OrganizationCode.Set {
		organizationCode = lo.ToPtr(params.OrganizationCode.Value)
	}

	out, err := a.AuthLogin.Execute(ctx, auth_usecase.AuthLoginInput{
		Provider:         string(provider),
		RedirectURL:      params.RedirectURL,
		RegistrationURL:  registrationURL,
		OrganizationCode: organizationCode,
	})
	if err != nil {
		return nil, err
	}

	headers := new(oas.AuthorizeFoundHeaders)
	headers.SetLocation(out.RedirectURL)
	return headers, nil
}

func (a *authHandler) DevAuthorize(ctx context.Context, params oas.DevAuthorizeParams) (oas.DevAuthorizeRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.DevAuthorize")
	defer span.End()

	var organizationCode *string
	if params.OrganizationCode.Set {
		organizationCode = lo.ToPtr(params.OrganizationCode.Value)
	}

	output, err := a.LoginForDev.Execute(ctx, auth_usecase.LoginForDevInput{
		Subject:          params.ID,
		OrganizationCode: organizationCode,
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

	output, err := a.AuthCallback.Execute(ctx, auth_usecase.CallbackInput{
		Provider: params.Provider,
		Code:     params.Code,
		State:    params.State,
	})
	if err != nil {
		return nil, err
	}

	headers := new(oas.OAuthCallbackFoundHeaders)
	headers.SetSetCookie(cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateSessionCookie(output.Token)}))
	headers.SetLocation(output.RedirectURL)
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
	_, err = a.Revoke.Execute(ctx, auth_usecase.RevokeInput{
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
	var orgType *int
	if claim.OrgType != nil {
		orgType = lo.ToPtr(*claim.OrgType)
	}

	return &oas.OAuthTokenInfoOK{
		Aud:                    claim.Audience(),
		Iat:                    claim.IssueAt().Format(time.RFC3339),
		Exp:                    claim.ExpiresAt().Format(time.RFC3339),
		Iss:                    claim.Issuer(),
		Sub:                    claim.Sub,
		Jti:                    sessID.String(),
		IsRegistered:           claim.IsRegistered,
		IsEmailVerified:        claim.IsEmailVerified,
		DisplayID:              utils.ToOpt[oas.OptString](claim.DisplayID),
		DisplayName:            utils.ToOpt[oas.OptString](claim.DisplayName),
		IconURL:                utils.ToOpt[oas.OptString](claim.IconURL),
		RequiredPasswordChange: claim.RequiredPasswordChange,
		OrgType:                utils.ToOptNil[oas.OptNilInt](orgType),
		OrganizationRole:       utils.ToOptNil[oas.OptNilString](claim.OrganizationRole),
		OrganizationCode:       utils.ToOptNil[oas.OptNilString](claim.OrganizationCode),
		OrganizationID:         utils.ToOptNil[oas.OptNilString](claim.OrganizationID),
	}, nil
}

// AuthAccountDetach 開発向け。退会処理を作るまでの代替。Subjectを付け替えることで、一度SSOしても再度SSOさせることができるやつ。
func (a *authHandler) AuthAccountDetach(ctx context.Context) (oas.AuthAccountDetachRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.AuthAccountDetach")
	defer span.End()

	claim := session.GetSession(ctx)
	sessID, err := claim.SessionID()
	if err != nil {
		return nil, messages.ForbiddenError
	}

	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}

	if err = a.DetachAccount.Execute(ctx, auth_usecase.DetachAccountInput{
		UserID: shared.UUID[user.User](userID),
	}); err != nil {
		return nil, err
	}

	// revoke
	_, err = a.Revoke.Execute(ctx, auth_usecase.RevokeInput{
		SessID: sessID,
	})
	if err != nil {
		return nil, err
	}

	res := new(oas.AuthAccountDetachOKHeaders)
	res.SetSetCookie(cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateRevokeCookie()}))
	return res, nil
}

// PasswordLogin
func (a *authHandler) PasswordLogin(ctx context.Context, req *oas.PasswordLoginReq) (oas.PasswordLoginRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.PasswordLogin")
	defer span.End()

	out, err := a.passwordLogin.Execute(ctx, auth_usecase.PasswordLoginInput{
		IDorEmail: req.IDOrEmail,
		Password:  req.Password,
	})
	if err != nil {
		return nil, err
	}

	res := http_utils.GetHTTPResponse(ctx)
	res.Header().Set("Set-Cookie", cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateSessionCookie(out.Token)})[0])
	return &oas.PasswordLoginOK{}, nil
}

// PasswordRegister
func (a *authHandler) PasswordRegister(ctx context.Context, req *oas.PasswordRegisterReq) (oas.PasswordRegisterRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.PasswordRegister")
	defer span.End()

	out, err := a.passwordRegister.Execute(ctx, auth_usecase.PasswordRegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	res := http_utils.GetHTTPResponse(ctx)
	res.Header().Set("Set-Cookie", cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateSessionCookie(out.Token)})[0])
	return &oas.PasswordRegisterOK{}, nil
}

// ChangePassword implements oas.AuthHandler.
func (a *authHandler) ChangePassword(ctx context.Context, params oas.ChangePasswordParams) (oas.ChangePasswordRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.ChangePassword")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}

	out, err := a.changePassword.Execute(ctx, auth_usecase.ChangePasswordInput{
		UserID:      shared.UUID[user.User](userID),
		OldPassword: params.OldPassword,
		NewPassword: params.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	if !out.Success {
		return nil, messages.InternalServerError
	}

	res := &oas.ChangePasswordOK{}
	return res, nil
}
