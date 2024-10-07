package handler

import (
	"context"
	"net/http"
	"net/url"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	"github.com/neko-dream/server/internal/presentation/oas"
	auth_usecase "github.com/neko-dream/server/internal/usecase/auth"
	http_utils "github.com/neko-dream/server/pkg/http"
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

// OAuthRevoke implements oas.AuthHandler.
func (a *authHandler) OAuthRevoke(ctx context.Context) (oas.OAuthRevokeRes, error) {
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

	w := http_utils.GetHTTPResponse(ctx)
	for _, c := range out.Cookies {
		http.SetCookie(w, c)
	}

	// 204 No Content
	return *new(oas.OAuthRevokeRes), nil
}

// OAuthTokenInfo implements oas.AuthHandler.
func (a *authHandler) OAuthTokenInfo(ctx context.Context) (*oas.OAuthTokenInfoOK, error) {
	claim := session.GetSession(ctx)
	sessID, err := claim.SessionID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if claim.IsExpired() {
		return nil, messages.TokenExpiredError
	}

	var (
		Aud      string = claim.Audience()
		Iat      string = time.NewTime(ctx, claim.IssueAt()).Format(ctx)
		Exp      string = time.NewTime(ctx, claim.ExpiresAt()).Format(ctx)
		Iss      string = claim.Issuer()
		Sub      string = claim.Sub
		Jti      string = sessID.String()
		IsVerify bool   = claim.IsVerify
	)

	return &oas.OAuthTokenInfoOK{
		Aud:         Aud,
		Iat:         Iat,
		Exp:         Exp,
		Iss:         Iss,
		Sub:         Sub,
		Jti:         Jti,
		IsVerify:    IsVerify,
		DisplayId:   utils.ToOptNil[oas.OptNilString](claim.DisplayID),
		DisplayName: utils.ToOptNil[oas.OptNilString](claim.DisplayName),
		IconURL:     utils.ToOptNil[oas.OptNilString](claim.IconURL),
	}, nil

}
