package organization_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type ListOrganizationAliasesInput struct {
	OrganizationID shared.UUID[organization.Organization]
}

type ListOrganizationAliasesOutput struct {
	Aliases []dto.OrganizationAlias
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

	aliasDTOs := make([]dto.OrganizationAlias, 0, len(aliases))
	for _, alias := range aliases {
		aliasDTOs = append(aliasDTOs, dto.OrganizationAlias{
			AliasID:   alias.AliasID().String(),
			AliasName: alias.AliasName(),
			CreatedAt: lo.ToPtr(alias.CreatedAt()),
		})
	}

	return &ListOrganizationAliasesOutput{
		Aliases: aliasDTOs,
	}, nil
}
