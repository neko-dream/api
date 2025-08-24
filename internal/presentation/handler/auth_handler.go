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
	"github.com/neko-dream/server/internal/domain/service"
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
	reactivate       auth_usecase.Reactivate

	authenticationService service.AuthenticationService
	authorizationService  service.AuthorizationService
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
	reactivate auth_usecase.Reactivate,

	authorizationService service.AuthorizationService,
	cookieManger cookie.CookieManager,
) oas.AuthHandler {
	return &authHandler{
		AuthLogin:            authLogin,
		AuthCallback:         authCallback,
		Revoke:               revoke,
		LoginForDev:          devLogin,
		DetachAccount:        detachAccount,
		authorizationService: authorizationService,
		CookieManager:        cookieManger,
		passwordLogin:        login,
		passwordRegister:     register,
		changePassword:       changePassword,
		reactivate:           reactivate,
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

// HandleAuthCallback implements oas.Handler.
func (a *authHandler) HandleAuthCallback(ctx context.Context, params oas.HandleAuthCallbackParams) (oas.HandleAuthCallbackRes, error) {
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

	headers := new(oas.HandleAuthCallbackFoundHeaders)
	headers.SetSetCookie(cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateSessionCookie(output.Token)}))
	headers.SetLocation(output.RedirectURL)
	return headers, nil
}

// RevokeToken ログアウト
func (a *authHandler) RevokeToken(ctx context.Context) (oas.RevokeTokenRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.OAuthTokenRevoke")
	defer span.End()

	authCtx, err := a.authorizationService.GetAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = a.Revoke.Execute(ctx, auth_usecase.RevokeInput{
		SessID: authCtx.SessionID,
	})
	if err != nil {
		return nil, err
	}

	headers := new(oas.RevokeTokenNoContent)
	headers.SetSetCookie(cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateRevokeCookie()}))
	return headers, nil
}

// GetTokenInfo implements oas.Handler.
func (a *authHandler) GetTokenInfo(ctx context.Context) (oas.GetTokenInfoRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.OAuthTokenInfo")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	if claim.IsExpired(ctx) {
		return nil, messages.TokenExpiredError
	}

	sessID, err := claim.SessionID()
	if err != nil {
		return nil, messages.InternalServerError
	}

	var orgType *int
	if claim.OrgType != nil {
		orgType = lo.ToPtr(*claim.OrgType)
	}

	return &oas.TokenClaim{
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

	authCtx, err := a.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID
	if err = a.DetachAccount.Execute(ctx, auth_usecase.DetachAccountInput{
		UserID: shared.UUID[user.User](userID),
	}); err != nil {
		return nil, err
	}

	// revoke
	_, err = a.Revoke.Execute(ctx, auth_usecase.RevokeInput{
		SessID: authCtx.SessionID,
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
		IDorEmail: req.IdOrEmail,
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

	authCtx, err := a.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID

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

func (a *authHandler) ReactivateUser(ctx context.Context) (oas.ReactivateUserRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "authHandler.ReactivateUser")
	defer span.End()

	// トークンから退会ユーザー情報を取得
	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	// 退会していないユーザーからのリクエストの場合
	if !claim.IsWithdrawn {
		return nil, messages.UserNotWithdrawn
	}

	// ユーザーIDを取得
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.InternalServerError
	}

	// 復活処理を実行
	output, err := a.reactivate.Execute(ctx, auth_usecase.ReactivateInput{
		UserID: userID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "reactivate.Execute")
		return nil, err
	}

	return &oas.ReactivateUserOK{
		Message: output.Message,
		User:    output.User.ToResponse(),
	}, nil
}
