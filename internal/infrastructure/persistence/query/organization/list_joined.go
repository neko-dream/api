package organization

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/application/query/organization_query"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type listJoinedOrganizationQuery struct {
	db *db.DBManager
}

func NewListJoinedOrganizationQuery(db *db.DBManager) organization_query.ListJoinedOrganizationQuery {
	return &listJoinedOrganizationQuery{db: db}
}

func (q *listJoinedOrganizationQuery) Execute(ctx context.Context, input organization_query.ListJoinedOrganizationInput) (*organization_query.ListJoinedOrganizationOutput, error) {
	ctx, span := otel.Tracer("organization").Start(ctx, "listJoinedOrganizationQuery.Execute")
	defer span.End()

	orgs, err := q.db.GetQueries(ctx).FindOrgUserByUserIDWithOrganization(ctx, input.UserID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &organization_query.ListJoinedOrganizationOutput{
				Organizations: []*dto.OrganizationResponse{},
			}, nil
		}
		utils.HandleError(ctx, err, "failed to find organization")
		return nil, err
	}

	orgRespList := make([]*dto.OrganizationResponse, 0, len(orgs))
	for _, org := range orgs {
		var orgResp dto.OrganizationResponse
		err := copier.CopyWithOption(&orgResp, org, copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
		})
		if err != nil {
			utils.HandleError(ctx, err, "failed to copy organization")
			return nil, err
		}

		orgRespList = append(orgRespList, &orgResp)
	}

	return &organization_query.ListJoinedOrganizationOutput{
		Organizations: orgRespList,
	}, nil
}
