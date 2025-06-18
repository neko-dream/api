package organization_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"go.opentelemetry.io/otel"
)

type ListOrganizationAliasesInput struct {
	OrganizationID shared.UUID[organization.Organization]
}

type OrganizationAliasDTO struct {
	AliasID   string
	AliasName string
	CreatedAt string
}

type ListOrganizationAliasesOutput struct {
	Aliases []OrganizationAliasDTO
}

type ListOrganizationAliasesUseCase struct {
	orgAliasService *service.OrganizationAliasService
}

func NewListOrganizationAliasesUseCase(
	orgAliasService *service.OrganizationAliasService,
) *ListOrganizationAliasesUseCase {
	return &ListOrganizationAliasesUseCase{
		orgAliasService: orgAliasService,
	}
}

func (u *ListOrganizationAliasesUseCase) Execute(
	ctx context.Context,
	input ListOrganizationAliasesInput,
) (*ListOrganizationAliasesOutput, error) {
	ctx, span := otel.Tracer("organization_usecase").Start(ctx, "ListOrganizationAliasesUseCase.Execute")
	defer span.End()

	aliases, err := u.orgAliasService.GetActiveAliases(ctx, input.OrganizationID)
	if err != nil {
		return nil, err
	}

	aliasDTOs := make([]OrganizationAliasDTO, 0, len(aliases))
	for _, alias := range aliases {
		aliasDTOs = append(aliasDTOs, OrganizationAliasDTO{
			AliasID:   alias.AliasID().String(),
			AliasName: alias.AliasName(),
			CreatedAt: alias.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &ListOrganizationAliasesOutput{
		Aliases: aliasDTOs,
	}, nil
}
