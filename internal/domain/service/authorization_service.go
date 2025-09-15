package service

import (
	"context"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/auth"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

// AuthorizationService ユーザー認可を扱うサービス
type AuthorizationService interface {
	// 認証不要でも呼び出し可能（認証されていない場合はnilを返す）
	GetAuthContext(ctx context.Context) (*auth.AuthenticationContext, error)

	// 認証していればUserIDが帰り、していなければnilが帰る
	GetUserID(ctx context.Context) (*shared.UUID[user.User], error)

	// 認証必須（認証されていない場合はエラー）
	RequireAuthentication(ctx context.Context) (*auth.AuthenticationContext, error)

	// 組織ロール必須
	RequireOrganizationRole(ctx context.Context, minRole organization.OrganizationUserRole) (*auth.AuthenticationContext, error)

	// スーパー管理者必須（RequireOrganizationRoleのエイリアス）
	RequireSuperAdmin(ctx context.Context) (*auth.AuthenticationContext, error)

	// オーナー権限必須（RequireOrganizationRoleのエイリアス）
	RequireOwner(ctx context.Context) (*auth.AuthenticationContext, error)

	// 管理者以上の権限必須（RequireOrganizationRoleのエイリアス）
	RequireAdmin(ctx context.Context) (*auth.AuthenticationContext, error)

	// 運営ユーザーかどうか
	IsKotohiro(ctx context.Context) bool

	// 認証されているかをチェック
	IsAuthenticated(ctx context.Context) bool

	// 組織コンテキスト内かをチェック
	IsInOrganization(ctx context.Context) bool
}

type authorizationService struct {
	authenticationService AuthenticationService
}

func NewAuthorizationService(
	authenticationService AuthenticationService,
) AuthorizationService {
	return &authorizationService{
		authenticationService: authenticationService,
	}
}

// GetAuthContext 認証不要でも呼び出し可能（認証されていない場合はnilを返す）
func (a *authorizationService) GetAuthContext(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.GetAuthContext")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, nil
	}

	return a.authenticationService.BuildAuthContext(ctx, claim)
}

func (a *authorizationService) GetUserID(ctx context.Context) (*shared.UUID[user.User], error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.GetUserID")
	defer span.End()

	authCtx, err := a.GetAuthContext(ctx)
	if err != nil {
		return nil, err
	}
	if authCtx == nil {
		return nil, nil
	}

	return &authCtx.UserID, nil
}

// RequireAuthentication 認証必須（認証されていない場合はエラー）
func (a *authorizationService) RequireAuthentication(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.RequireAuthentication")
	defer span.End()

	claim, err := a.GetAuthContext(ctx)
	if claim == nil || err != nil {
		return nil, messages.ForbiddenError
	}

	return claim, nil
}

// RequireOrganizationRole 組織ロール必須
func (a *authorizationService) RequireOrganizationRole(ctx context.Context, minRole organization.OrganizationUserRole) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.RequireOrganizationRole")
	defer span.End()

	authCtx, err := a.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	if !authCtx.IsInOrganization() {
		return nil, messages.OrganizationContextRequired
	}

	if !authCtx.HasOrganizationRole(minRole) {
		return nil, messages.OrganizationPermissionDenied
	}

	return authCtx, nil
}

// RequireSuperAdmin スーパー管理者必須（RequireOrgRoleのエイリアス）
func (a *authorizationService) RequireSuperAdmin(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.RequireSuperAdmin")
	defer span.End()

	return a.RequireOrganizationRole(ctx, organization.OrganizationUserRoleSuperAdmin)
}

// RequireOwner オーナー権限必須（RequireOrgRoleのエイリアス）
func (a *authorizationService) RequireOwner(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.RequireOwner")
	defer span.End()

	return a.RequireOrganizationRole(ctx, organization.OrganizationUserRoleOwner)
}

// RequireAdmin 管理者以上の権限必須（RequireOrgRoleのエイリアス）
func (a *authorizationService) RequireAdmin(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.RequireAdmin")
	defer span.End()

	return a.RequireOrganizationRole(ctx, organization.OrganizationUserRoleAdmin)
}

// IsAuthenticated 認証されているかをチェック
func (a *authorizationService) IsAuthenticated(ctx context.Context) bool {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.IsAuthenticated")
	defer span.End()

	_, err := a.GetAuthContext(ctx)
	return err == nil
}

// IsInOrganization 組織コンテキスト内かをチェック
func (a *authorizationService) IsInOrganization(ctx context.Context) bool {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.IsInOrganization")
	defer span.End()

	authCtx, err := a.GetAuthContext(ctx)
	if err != nil {
		return false
	}
	return authCtx.IsInOrganization()
}

// IsKotohio implements AuthorizationService.
func (a *authorizationService) IsKotohiro(ctx context.Context) bool {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.IsKotohiro")
	defer span.End()

	authCtx, err := a.GetAuthContext(ctx)
	if err != nil {
		return false
	}

	if !authCtx.IsInOrganization() {
		return false
	}

	if !authCtx.IsKotohiro() {
		return false
	}

	return true
}
