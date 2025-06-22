package service

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type authenticationServiceImpl struct{}

func NewAuthenticationService() service.AuthenticationService {
	return &authenticationServiceImpl{}
}

func (a *authenticationServiceImpl) GetCurrentUser(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authenticationServiceImpl.GetCurrentUser")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, service.ErrNotAuthenticated
	}

	return a.claimToAuthenticationContext(ctx, claim)
}

func (a *authenticationServiceImpl) RequireAuthentication(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authenticationServiceImpl.RequireAuthentication")
	defer span.End()

	return a.GetCurrentUser(ctx)
}

func (a *authenticationServiceImpl) RequireOrganizationRole(ctx context.Context, minRole organization.OrganizationUserRole) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authenticationServiceImpl.RequireOrganizationRole")
	defer span.End()

	authCtx, err := a.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	if !authCtx.IsInOrganization() {
		return nil, service.ErrNotInOrganization
	}

	if !authCtx.HasOrganizationRole(minRole) {
		return nil, service.ErrInsufficientPermissions
	}

	return authCtx, nil
}

func (a *authenticationServiceImpl) RequireSuperAdmin(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authenticationServiceImpl.RequireSuperAdmin")
	defer span.End()

	return a.RequireOrganizationRole(ctx, organization.OrganizationUserRoleSuperAdmin)
}

func (a *authenticationServiceImpl) RequireOwner(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authenticationServiceImpl.RequireOwner")
	defer span.End()

	return a.RequireOrganizationRole(ctx, organization.OrganizationUserRoleOwner)
}

func (a *authenticationServiceImpl) RequireAdmin(ctx context.Context) (*auth.AuthenticationContext, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authenticationServiceImpl.RequireAdmin")
	defer span.End()

	return a.RequireOrganizationRole(ctx, organization.OrganizationUserRoleAdmin)
}

func (a *authenticationServiceImpl) IsAuthenticated(ctx context.Context) bool {
	ctx, span := otel.Tracer("service").Start(ctx, "authenticationServiceImpl.IsAuthenticated")
	defer span.End()

	_, err := a.GetCurrentUser(ctx)
	return err == nil
}

func (a *authenticationServiceImpl) IsInOrganization(ctx context.Context) bool {
	ctx, span := otel.Tracer("service").Start(ctx, "authenticationServiceImpl.IsInOrganization")
	defer span.End()

	authCtx, err := a.GetCurrentUser(ctx)
	if err != nil {
		return false
	}
	return authCtx.IsInOrganization()
}

// claimToAuthenticationContext converts a session.Claim to auth.AuthenticationContext
func (a *authenticationServiceImpl) claimToAuthenticationContext(ctx context.Context, claim *session.Claim) (*auth.AuthenticationContext, error) {
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, service.ErrNotAuthenticated
	}

	sessionID, err := claim.SessionID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.SessionID")
		return nil, service.ErrNotAuthenticated
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
