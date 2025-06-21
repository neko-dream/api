package handler

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
)

// 認証コンテキストを取得
func getAuthenticationContext(authService service.AuthenticationService, ctx context.Context) (*auth.AuthenticationContext, error) {
	authCtx, err := authService.GetCurrentUser(ctx)
	if err != nil {
		if err == service.ErrNotAuthenticated {
			return nil, messages.ForbiddenError
		}
		return nil, err
	}
	return authCtx, nil
}

// 認証されていることを確認
func requireAuthentication(authService service.AuthenticationService, ctx context.Context) (*auth.AuthenticationContext, error) {
	authCtx, err := authService.RequireAuthentication(ctx)
	if err != nil {
		if err == service.ErrNotAuthenticated {
			return nil, messages.ForbiddenError
		}
		return nil, err
	}
	return authCtx, nil
}

// 指定された組織役割以上の権限を要求
func requireOrganizationRole(authService service.AuthenticationService, ctx context.Context, minRole shared.OrganizationUserRole) (*auth.AuthenticationContext, error) {
	authCtx, err := authService.RequireOrganizationRole(ctx, minRole)
	if err != nil {
		if err == service.ErrNotAuthenticated {
			return nil, messages.ForbiddenError
		}
		if err == service.ErrInsufficientPermissions {
			return nil, messages.ForbiddenError
		}
		if err == service.ErrNotInOrganization {
			return nil, messages.ForbiddenError
		}
		return nil, err
	}
	return authCtx, nil
}
