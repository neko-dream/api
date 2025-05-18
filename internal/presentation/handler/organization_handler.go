package handler

import (
	"context"

	"github.com/neko-dream/server/internal/application/command/organization_command"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

type organizationHandler struct {
	create organization_command.CreateOrganizationCommand
	invite organization_command.InviteOrganizationCommand
	add    organization_command.InviteOrganizationForUserCommand
}

func NewOrganizationHandler(
	create organization_command.CreateOrganizationCommand,
	invite organization_command.InviteOrganizationCommand,
	add organization_command.InviteOrganizationForUserCommand,
) oas.OrganizationHandler {
	return &organizationHandler{
		create: create,
		invite: invite,
		add:    add,
	}
}

// CreateOrganizations implements oas.OrganizationHandler.
func (o *organizationHandler) CreateOrganizations(ctx context.Context, req oas.OptCreateOrganizationsReq) (oas.CreateOrganizationsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.CreateOrganizations")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	_, err = o.create.Execute(ctx, organization_command.CreateOrganizationInput{
		UserID: userID,
		Name:   req.Value.Name,
		Type:   req.Value.OrgType.Value,
	})
	if err != nil {
		return nil, err
	}

	res := &oas.CreateOrganizationsOK{}
	return res, nil
}

// InviteOrganization implements oas.OrganizationHandler.
func (o *organizationHandler) InviteOrganization(ctx context.Context, req oas.OptInviteOrganizationReq, params oas.InviteOrganizationParams) (oas.InviteOrganizationRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.InviteOrganization")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	organizationID, err := shared.ParseUUID[organization.Organization](params.OrganizationID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	_, err = o.invite.Execute(ctx, organization_command.InviteOrganizationInput{
		UserID:         userID,
		OrganizationID: organizationID,
		Email:          req.Value.Email,
		Role:           req.Value.Role,
	})
	if err != nil {
		return nil, err
	}

	res := &oas.InviteOrganizationOK{}
	return res, nil
}

// InviteOrganizationForUser implements oas.OrganizationHandler.
func (o *organizationHandler) InviteOrganizationForUser(ctx context.Context, req oas.OptInviteOrganizationForUserReq, params oas.InviteOrganizationForUserParams) (oas.InviteOrganizationForUserRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.InviteOrganizationForUser")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	organizationID, err := shared.ParseUUID[organization.Organization](params.OrganizationID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	_, err = o.add.Execute(ctx, organization_command.InviteOrganizationForUserInput{
		UserID:         userID,
		OrganizationID: organizationID,
		DisplayID:      req.Value.DisplayID,
		Role:           int(req.Value.Role),
	})
	if err != nil {
		return nil, err
	}

	res := &oas.InviteOrganizationForUserOK{
		Success: true,
	}
	return res, nil
}

// GetOrganizations 所属組織一覧
func (o *organizationHandler) GetOrganizations(ctx context.Context) (oas.GetOrganizationsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.GetOrganizations")
	defer span.End()

	_ = ctx
	return nil, &messages.APIError{
		StatusCode: 501,
		Code:       "ORG-0001",
		Message:    "この機能はまだ実装されていません",
	}
}
