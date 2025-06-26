package organization_query

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

type ListOrganizationUsersQuery interface {
	Execute(context.Context, ListOrganizationUsersInput) (*ListOrganizationUsersOutput, error)
}

type ListOrganizationUsersInput struct {
	OrganizationID shared.UUID[organization.Organization]
}

type ListOrganizationUsersOutput struct {
	Users []oas.OrganizationUser
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

	var users []oas.OrganizationUser
	for _, row := range result {
		iconURL := oas.OptNilString{}
		if row.IconUrl.Valid {
			iconURL = oas.NewOptNilString(row.IconUrl.String)
		}

		displayID := ""
		if row.DisplayID.Valid {
			displayID = row.DisplayID.String
		}

		displayName := ""
		if row.DisplayName.Valid {
			displayName = row.DisplayName.String
		}

		users = append(users, oas.OrganizationUser{
			UserID:      row.UserID.String(),
			DisplayID:   displayID,
			DisplayName: displayName,
			IconURL:     iconURL,
			Role:        int(row.Role),
			RoleName:    organization.RoleToName(organization.OrganizationUserRole(row.Role)),
		})
	}

	return &ListOrganizationUsersOutput{
		Users: users,
	}, nil
}