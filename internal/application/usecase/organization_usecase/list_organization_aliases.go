package organization_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"go.opentelemetry.io/otel"
)

// ListOrganizationAliasesInput 入力
type ListOrganizationAliasesInput struct {
	OrganizationID string
}

// OrganizationAliasDTO エイリアスDTO
type OrganizationAliasDTO struct {
	AliasID   string
	AliasName string
	CreatedAt string
}

// ListOrganizationAliasesOutput 出力
type ListOrganizationAliasesOutput struct {
	Aliases []OrganizationAliasDTO
}

// ListOrganizationAliasesUseCase エイリアス一覧取得ユースケース
type ListOrganizationAliasesUseCase struct {
	orgAliasService *service.OrganizationAliasService
}

// NewListOrganizationAliasesUseCase コンストラクタ
func NewListOrganizationAliasesUseCase(
	orgAliasService *service.OrganizationAliasService,
) *ListOrganizationAliasesUseCase {
	return &ListOrganizationAliasesUseCase{
		orgAliasService: orgAliasService,
	}
}

// Execute エイリアス一覧取得を実行
func (u *ListOrganizationAliasesUseCase) Execute(
	ctx context.Context,
	input ListOrganizationAliasesInput,
) (*ListOrganizationAliasesOutput, error) {
	ctx, span := otel.Tracer("organization_usecase").Start(ctx, "ListOrganizationAliasesUseCase.Execute")
	defer span.End()

	// 組織ID解析
	orgID, err := shared.ParseUUID[organization.Organization](input.OrganizationID)
	if err != nil {
		return nil, err
	}

	// アクティブなエイリアス一覧を取得
	aliases, err := u.orgAliasService.GetActiveAliases(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// DTOに変換
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
