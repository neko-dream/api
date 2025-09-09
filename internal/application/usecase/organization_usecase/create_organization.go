package organization_usecase

import (
	"context"
	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	organization_svc "github.com/neko-dream/server/internal/domain/service/organization"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type CreateOrganizationCommand interface {
	Execute(ctx context.Context, input CreateOrganizationInput) (*CreateOrganizationOutput, error)
}

type CreateOrganizationInput struct {
	UserID    shared.UUID[user.User]
	Name      string
	Code      string
	Type      int
	IconImage *multipart.FileHeader
}

type CreateOrganizationOutput struct {
}

type createOrganizationInteractor struct {
	organizationService organization_svc.OrganizationService
}

func NewCreateOrganizationInteractor(
	organizationService organization_svc.OrganizationService,
) CreateOrganizationCommand {
	return &createOrganizationInteractor{
		organizationService: organizationService,
	}
}

func (c *createOrganizationInteractor) Execute(ctx context.Context, input CreateOrganizationInput) (*CreateOrganizationOutput, error) {
	ctx, span := otel.Tracer("organization_command").Start(ctx, "createOrganizationInteractor.Execute")
	defer span.End()

	orgType := organization.OrganizationType(input.Type)
	ownerID := input.UserID

	_, err := c.organizationService.CreateOrganization(ctx, input.Name, input.Code, input.IconImage, orgType, ownerID)
	if err != nil {
		utils.HandleError(ctx, err, "CreateOrganization")
		return nil, err
	}

	return &CreateOrganizationOutput{}, nil
}
