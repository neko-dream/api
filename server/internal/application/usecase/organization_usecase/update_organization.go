package organization_usecase

import (
	"context"
	"mime/multipart"

	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	organization_svc "github.com/neko-dream/api/internal/domain/service/organization"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type UpdateOrganizationCommand interface {
	Execute(ctx context.Context, input UpdateOrganizationInput) error
}

type UpdateOrganizationInput struct {
	UserID         shared.UUID[user.User]
	OrganizationID shared.UUID[organization.Organization]
	Name           string
	IconImage      *multipart.FileHeader
}

type UpdateOrganizationInteractor struct {
	organizationService organization_svc.OrganizationService
}

func NewUpdateOrganizationInteractor(
	organizationService organization_svc.OrganizationService,
) UpdateOrganizationCommand {
	return &UpdateOrganizationInteractor{
		organizationService: organizationService,
	}
}

func (c *UpdateOrganizationInteractor) Execute(ctx context.Context, input UpdateOrganizationInput) error {
	ctx, span := otel.Tracer("organization_command").Start(ctx, "UpdateOrganizationInteractor.Execute")
	defer span.End()

	err := c.organizationService.UpdateOrganization(ctx, input.OrganizationID, input.Name, input.IconImage)
	if err != nil {
		utils.HandleError(ctx, err, "UpdateOrganization")
		return err
	}

	return nil
}
