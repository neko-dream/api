package handler

import (
	"context"
	"net/http"

	"github.com/neko-dream/api/internal/application/usecase/admin_usecase"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/infrastructure/http/cookie"
	"github.com/neko-dream/api/internal/presentation/oas"
	cookie_utils "github.com/neko-dream/api/pkg/cookie"
	http_utils "github.com/neko-dream/api/pkg/http"
	"go.opentelemetry.io/otel"
)

type adminHandler struct {
	adminLogin    admin_usecase.AdminLogin
	adminRefresh  admin_usecase.AdminRefresh
	adminLogout   admin_usecase.AdminLogout
	cookieManager cookie.CookieManager
}

func NewAdminHandler(
	adminLogin admin_usecase.AdminLogin,
	adminRefresh admin_usecase.AdminRefresh,
	adminLogout admin_usecase.AdminLogout,
	cookieManager cookie.CookieManager,
) oas.AdminHandler {
	return &adminHandler{
		adminLogin:    adminLogin,
		adminRefresh:  adminRefresh,
		adminLogout:   adminLogout,
		cookieManager: cookieManager,
	}
}

// sessionCookies は access(AdminSessionId) と refresh(AdminRefreshToken) の Set-Cookie を生成する。
func (h *adminHandler) sessionCookies(out *admin_usecase.AdminAuthOutput) []string {
	return cookie_utils.EncodeCookies([]*http.Cookie{
		h.cookieManager.CreateAdminSessionCookie(out.AccessToken),
		h.cookieManager.CreateAdminRefreshCookie(out.RefreshToken),
	})
}

func (h *adminHandler) tokenInfo(out *admin_usecase.AdminAuthOutput) oas.AdminTokenInfo {
	return oas.AdminTokenInfo{
		Sub:         out.AdminUser.AdminUserID.String(),
		Email:       out.AdminUser.Email,
		DisplayName: out.AdminUser.DisplayName,
		Role:        int(out.AdminUser.Role),
		IsAdmin:     true,
	}
}

// adminRefreshTokenFromCtx は HTTP リクエストから AdminRefreshToken Cookie を読む。
func adminRefreshTokenFromCtx(ctx context.Context) string {
	req := http_utils.GetHTTPRequest(ctx)
	if req == nil {
		return ""
	}
	c, err := req.Cookie(cookie.AdminRefreshCookieName)
	if err != nil {
		return ""
	}
	return c.Value
}

// AdminLogin implements adminLogin operation.
func (h *adminHandler) AdminLogin(ctx context.Context, req *oas.AdminLoginReq) (oas.AdminLoginRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "adminHandler.AdminLogin")
	defer span.End()

	out, err := h.adminLogin.Execute(ctx, admin_usecase.AdminLoginInput{
		IDToken: req.IdToken,
	})
	if err != nil {
		return &oas.AdminLoginUnauthorized{}, nil
	}

	return &oas.AdminTokenInfoHeaders{
		SetCookie: h.sessionCookies(out),
		Response:  h.tokenInfo(out),
	}, nil
}

// AdminRefresh implements adminRefresh operation.
func (h *adminHandler) AdminRefresh(ctx context.Context) (oas.AdminRefreshRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "adminHandler.AdminRefresh")
	defer span.End()

	out, err := h.adminRefresh.Execute(ctx, adminRefreshTokenFromCtx(ctx))
	if err != nil {
		return &oas.AdminRefreshUnauthorized{}, nil
	}

	return &oas.AdminTokenInfoHeaders{
		SetCookie: h.sessionCookies(out),
		Response:  h.tokenInfo(out),
	}, nil
}

// GetAdminTokenInfo implements getAdminTokenInfo operation.
func (h *adminHandler) GetAdminTokenInfo(ctx context.Context) (oas.GetAdminTokenInfoRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "adminHandler.GetAdminTokenInfo")
	defer span.End()

	claim := admin.GetAdminSession(ctx)
	if claim == nil {
		return &oas.GetAdminTokenInfoForbidden{}, nil
	}
	return &oas.AdminTokenInfo{
		Sub:         claim.Sub,
		Email:       claim.Email,
		DisplayName: claim.DisplayName,
		Role:        claim.Role,
		IsAdmin:     true,
	}, nil
}

// AdminLogout implements adminLogout operation.
func (h *adminHandler) AdminLogout(ctx context.Context) (*oas.AdminLogoutNoContent, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "adminHandler.AdminLogout")
	defer span.End()

	// refresh セッションをサーバー側で失効させる(エラーでも cookie は消す)。
	_ = h.adminLogout.Execute(ctx, adminRefreshTokenFromCtx(ctx))

	return &oas.AdminLogoutNoContent{
		SetCookie: cookie_utils.EncodeCookies([]*http.Cookie{
			h.cookieManager.CreateAdminRevokeCookie(),
			h.cookieManager.CreateAdminRefreshRevokeCookie(),
		}),
	}, nil
}
