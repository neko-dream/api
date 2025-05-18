package organization_query

import (
	"context"

	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type ListJoinedOrganizationQuery interface {
	Execute(context.Context, ListJoinedOrganizationInput) (*ListJoinedOrganizationOutput, error)
}

type ListJoinedOrganizationInput struct {
	UserID shared.UUID[user.User]
}

type ListJoinedOrganizationOutput struct {
	Organizations []*dto.OrganizationResponse
}
