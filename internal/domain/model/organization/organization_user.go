package organization

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type OrganizationUserRepository interface {
	// OrganizationUserの取得
	FindByOrganizationIDAndUserID(ctx context.Context, orgID shared.UUID[Organization], userID shared.UUID[user.User]) (*OrganizationUser, error)
	FindByOrganizationID(ctx context.Context, orgID shared.UUID[Organization]) ([]*OrganizationUser, error)
	FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]*OrganizationUser, error)

	// OrganizationUserの作成・更新・削除
	Create(ctx context.Context, orgUser *OrganizationUser) error
}

type OrganizationUserRole int

const (
	OrganizationUserRoleMember OrganizationUserRole = iota + 1
	OrganizationUserRoleAdmin
	OrganizationUserRoleOwner
	OrganizationUserRoleSuperAdmin
)

type OrganizationUser struct {
	OrganizationID shared.UUID[Organization]
	UserID         shared.UUID[user.User]
	Role           OrganizationUserRole
}
