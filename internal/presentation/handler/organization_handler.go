package handler

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/command/organization_command"
	"go.opentelemetry.io/otel"
)

type organizationHandler struct {
	create organization_command.CreateOrganizationCommand
	invite organization_command.InviteOrganizationCommand
}

func NewOrganizationHandler(
	create organization_command.CreateOrganizationCommand,
	invite organization_command.InviteOrganizationCommand,
) oas.OrganizationHandler {
	return &organizationHandler{
		create: create,
		invite: invite,
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
