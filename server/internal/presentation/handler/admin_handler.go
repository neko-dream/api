package handler

import (
	"context"
	"errors"
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
	adminLogin      admin_usecase.AdminLogin
	adminRefresh    admin_usecase.AdminRefresh
	adminLogout     admin_usecase.AdminLogout
	listAdminUsers  admin_usecase.ListAdminUsers
	createAdminUser admin_usecase.CreateAdminUser
	cookieManager   cookie.CookieManager
	adminUserRepo   admin.AdminUserRepository
}

func NewAdminHandler(
	adminLogin admin_usecase.AdminLogin,
	adminRefresh admin_usecase.AdminRefresh,
	adminLogout admin_usecase.AdminLogout,
	listAdminUsers admin_usecase.ListAdminUsers,
	createAdminUser admin_usecase.CreateAdminUser,
	cookieManager cookie.CookieManager,
	adminUserRepo admin.AdminUserRepository,
) oas.AdminHandler {
	return &adminHandler{
		adminLogin:      adminLogin,
		adminRefresh:    adminRefresh,
		adminLogout:     adminLogout,
		listAdminUsers:  listAdminUsers,
		createAdminUser: createAdminUser,
		cookieManager:   cookieManager,
		adminUserRepo:   adminUserRepo,
	}
}

// ListAdminUsers implements listAdminUsers operation.
func (h *adminHandler) ListAdminUsers(ctx context.Context) (oas.ListAdminUsersRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "adminHandler.ListAdminUsers")
	defer span.End()

	users, err := h.listAdminUsers.Execute(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]oas.AdminUserItem, 0, len(users))
	for _, u := range users {
		items = append(items, oas.AdminUserItem{
			ID:          u.AdminUserID.String(),
			Email:       u.Email,
			DisplayName: u.DisplayName,
		})
	}
	return &oas.AdminUserListResponse{AdminUsers: items}, nil
}

// CreateAdminUser implements createAdminUser operation.
func (h *adminHandler) CreateAdminUser(ctx context.Context, req *oas.CreateAdminUserRequest) (oas.CreateAdminUserRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "adminHandler.CreateAdminUser")
	defer span.End()

	out, err := h.createAdminUser.Execute(ctx, admin_usecase.CreateAdminUserInput{
		Email:       req.Email,
		DisplayName: req.DisplayName.Or(""),
	})
	if err != nil {
		if errors.Is(err, admin_usecase.ErrAdminUserAlreadyExists) {
			return &oas.CreateAdminUserConflict{}, nil
		}
		return nil, err
	}
	return &oas.AdminUserItem{
		ID:          out.AdminUserID.String(),
		Email:       out.Email,
		DisplayName: out.DisplayName,
	}, nil
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
		if errors.Is(err, admin_usecase.ErrInvalidAdminCredentials) {
			return &oas.AdminLoginUnauthorized{}, nil
		}
		return nil, err
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
		if errors.Is(err, admin_usecase.ErrInvalidRefreshToken) {
			return &oas.AdminRefreshUnauthorized{}, nil
		}
		return nil, err
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
	id, err := claim.AdminUserID()
	if err != nil {
		return &oas.GetAdminTokenInfoForbidden{}, nil
	}
	user, err := h.adminUserRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &oas.GetAdminTokenInfoForbidden{}, nil
	}
	return &oas.AdminTokenInfo{
		Sub:         claim.Sub,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		IsAdmin:     true,
	}, nil
}

// AdminLogout implements adminLogout operation.
func (h *adminHandler) AdminLogout(ctx context.Context) (*oas.AdminLogoutNoContent, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "adminHandler.AdminLogout")
	defer span.End()

	// Do not report success if server-side refresh revocation failed.
	if err := h.adminLogout.Execute(ctx, adminRefreshTokenFromCtx(ctx)); err != nil {
		return nil, err
	}

	return &oas.AdminLogoutNoContent{
		SetCookie: cookie_utils.EncodeCookies([]*http.Cookie{
			h.cookieManager.CreateAdminRevokeCookie(),
			h.cookieManager.CreateAdminRefreshRevokeCookie(),
		}),
	}, nil
}
