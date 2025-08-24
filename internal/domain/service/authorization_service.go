package service

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

// AuthorizationService ユーザー認可を扱うサービス
type AuthorizationService interface {
	// 認証不要でも呼び出し可能（認証されていない場合はnilを返す）
	GetAuthContext(ctx context.Context) (*auth.AuthenticationContext, error)

	// 認証必須（認証されていない場合はエラー）
	RequireAuth(ctx context.Context) (*auth.AuthenticationContext, error)

	// 組織ロール必須
	RequireOrgRole(ctx context.Context, minRole organization.OrganizationUserRole) (*auth.AuthenticationContext, error)

	// スーパー管理者必須（RequireOrgRoleのエイリアス）
	RequireSuperAdmin(ctx context.Context) (*auth.AuthenticationContext, error)

	// オーナー権限必須（RequireOrgRoleのエイリアス）
	RequireOwner(ctx context.Context) (*auth.AuthenticationContext, error)

	// 管理者以上の権限必須（RequireOrgRoleのエイリアス）
	RequireAdmin(ctx context.Context) (*auth.AuthenticationContext, error)

	// 認証されているかをチェック
	IsAuthenticated(ctx context.Context) bool

	// 組織コンテキスト内かをチェック
	IsInOrganization(ctx context.Context) bool
}
type authorizationService struct{}

func NewAuthorizationService() AuthorizationService {
	return &authorizationService{}
}

// GetAuthContext 認証不要でも呼び出し可能（認証されていない場合はnilを返す）
func (a *authorizationService) GetAuthContext(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.GetAuthContext")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	return a.claimToAuthenticationContext(ctx, claim)
}

// RequireAuth 認証必須（認証されていない場合はエラー）
func (a *authorizationService) RequireAuth(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.RequireAuth")
	defer span.End()

	return a.GetAuthContext(ctx)
}

// RequireOrgRole 組織ロール必須
func (a *authorizationService) RequireOrgRole(ctx context.Context, minRole organization.OrganizationUserRole) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.RequireOrgRole")
	defer span.End()

	authCtx, err := a.RequireAuth(ctx)
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

	return a.RequireOrgRole(ctx, organization.OrganizationUserRoleSuperAdmin)
}

// RequireOwner オーナー権限必須（RequireOrgRoleのエイリアス）
func (a *authorizationService) RequireOwner(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.RequireOwner")
	defer span.End()

	return a.RequireOrgRole(ctx, organization.OrganizationUserRoleOwner)
}

// RequireAdmin 管理者以上の権限必須（RequireOrgRoleのエイリアス）
func (a *authorizationService) RequireAdmin(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authorizationService.RequireAdmin")
	defer span.End()

	return a.RequireOrgRole(ctx, organization.OrganizationUserRoleAdmin)
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

// claimToAuthenticationContext converts a session.Claim to auth.AuthenticationContext
func (a *authorizationService) claimToAuthenticationContext(ctx context.Context, claim *session.Claim) (*auth.AuthenticationContext, error) {
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	sessionID, err := claim.SessionID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.SessionID")
		return nil, messages.ForbiddenError
	}

	authCtx := &auth.AuthenticationContext{
		UserID:                 userID,
		SessionID:              sessionID,
		DisplayName:            claim.DisplayName,
		DisplayID:              claim.DisplayID,
		IconURL:                claim.IconURL,
		IsRegistered:           claim.IsRegistered,
		IsEmailVerified:        claim.IsEmailVerified,
		RequiredPasswordChange: claim.RequiredPasswordChange,
	}

	// 組織コンテキストがある場合は設定
	if claim.OrganizationID != nil {
		authCtx.OrganizationID = lo.ToPtr(shared.MustParseUUID[organization.Organization](*claim.OrganizationID))
	}

	if claim.OrganizationCode != nil {
		authCtx.OrganizationCode = claim.OrganizationCode
	}

	if claim.OrganizationRole != nil {
		role := organization.NameToRole(*claim.OrganizationRole)
		authCtx.OrganizationRole = &role
	}

	return authCtx, nil
}
