package organization_query

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type ListOrganizationUsersQuery interface {
	Execute(context.Context, ListOrganizationUsersInput) (*ListOrganizationUsersOutput, error)
}

type ListOrganizationUsersInput struct {
	OrganizationID shared.UUID[organization.Organization]
}

type ListOrganizationUsersOutput struct {
	Organizations []dto.OrganizationResponse
}

type listOrganizationUsersQuery struct {
	db *db.DBManager
}

func NewListOrganizationUsersQuery(db *db.DBManager) ListOrganizationUsersQuery {
	return &listOrganizationUsersQuery{
		db: db,
	}
}

func (q *listOrganizationUsersQuery) Execute(ctx context.Context, input ListOrganizationUsersInput) (*ListOrganizationUsersOutput, error) {
	ctx, span := otel.Tracer("query").Start(ctx, "listOrganizationUsersQuery.Execute")
	defer span.End()

	// 組織のユーザー一覧を詳細情報付きで取得
	result, err := q.db.GetQueries(ctx).FindOrganizationUsersWithDetails(ctx, input.OrganizationID.UUID())
	if err != nil {
		return nil, err
	}

	var orgRespList []dto.OrganizationResponse
	for _, org := range result {
		var orgResp dto.OrganizationResponse
		err := copier.CopyWithOption(&orgResp, org, copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
		})
		if err != nil {
			utils.HandleError(ctx, err, "failed to copy organization")
			return nil, err
		}
		orgResp.OrganizationUser.SetRoleName(int(org.OrganizationUser.Role))
		orgRespList = append(orgRespList, orgResp)
	}

	return &ListOrganizationUsersOutput{
		Organizations: orgRespList,
	}, nil
}
