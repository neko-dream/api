package handler

import (
	"context"
	"net/http"

	"github.com/neko-dream/api/internal/application/usecase/admin_usecase"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/infrastructure/http/cookie"
	"github.com/neko-dream/api/internal/presentation/oas"
	cookie_utils "github.com/neko-dream/api/pkg/cookie"
	"go.opentelemetry.io/otel"
)

type adminHandler struct {
	adminLogin    admin_usecase.AdminLogin
	cookieManager cookie.CookieManager
}

func NewAdminHandler(
	adminLogin admin_usecase.AdminLogin,
	cookieManager cookie.CookieManager,
) oas.AdminHandler {
	return &adminHandler{
		adminLogin:    adminLogin,
		cookieManager: cookieManager,
	}
}

// AdminLogin implements adminLogin operation.
func (h *adminHandler) AdminLogin(ctx context.Context, req *oas.AdminLoginReq) (oas.AdminLoginRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "adminHandler.AdminLogin")
	defer span.End()

	out, err := h.adminLogin.Execute(ctx, admin_usecase.AdminLoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		// 認証失敗は 401（資格情報の有無を区別しない）
		return &oas.AdminLoginUnauthorized{}, nil
	}

	sessionCookie := h.cookieManager.CreateAdminSessionCookie(out.Token)
	return &oas.AdminTokenInfoHeaders{
		SetCookie: cookie_utils.EncodeCookies([]*http.Cookie{sessionCookie}),
		Response: oas.AdminTokenInfo{
			Sub:         out.AdminUser.AdminUserID.String(),
			Email:       out.AdminUser.Email,
			DisplayName: out.AdminUser.DisplayName,
			Role:        int(out.AdminUser.Role),
			IsAdmin:     true,
		},
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

	revokeCookie := h.cookieManager.CreateAdminRevokeCookie()
	return &oas.AdminLogoutNoContent{
		SetCookie: cookie_utils.EncodeCookies([]*http.Cookie{revokeCookie}),
	}, nil
}
