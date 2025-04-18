package organization

import (
	"context"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type OrganizationUserRepository interface {
	// OrganizationUserの取得
	FindByOrganizationIDAndUserID(ctx context.Context, orgID shared.UUID[Organization], userID shared.UUID[user.User]) (*OrganizationUser, error)
	FindByOrganizationID(ctx context.Context, orgID shared.UUID[Organization]) ([]*OrganizationUser, error)
	FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]*OrganizationUser, error)

	// OrganizationUserの作成・更新・削除
	Create(ctx context.Context, orgUser OrganizationUser) error
}

type OrganizationUserRole int

const (
	OrganizationUserRoleMember OrganizationUserRole = iota + 1
	OrganizationUserRoleAdmin
	OrganizationUserRoleOwner
	OrganizationUserRoleSuperAdmin
)

type OrganizationUser struct {
	OrganizationUserID shared.UUID[OrganizationUser]
	OrganizationID     shared.UUID[Organization]
	UserID             shared.UUID[user.User]
	Role               OrganizationUserRole
}

// NewOrganizationUser は新しいOrganizationUserを作成するのじゃ
func NewOrganizationUser(
	organizationUserID shared.UUID[OrganizationUser],
	organizationID shared.UUID[Organization],
	userID shared.UUID[user.User],
	role OrganizationUserRole,
) *OrganizationUser {
	return &OrganizationUser{
		OrganizationUserID: organizationUserID,
		OrganizationID:     organizationID,
		UserID:             userID,
		Role:               role,
	}
}

// SetRole はユーザーのロールを設定するのじゃ
func (ou *OrganizationUser) SetRole(role OrganizationUserRole) error {
	if role < OrganizationUserRoleMember || role > OrganizationUserRoleSuperAdmin {
		return errors.New("invalid role")
	}
	ou.Role = role
	return nil
}

// HasPermissionToChangeRoleTo は現在のユーザーが指定されたロールに変更する権限があるかをチェックするのじゃ
func (ou *OrganizationUser) HasPermissionToChangeRoleTo(targetRole OrganizationUserRole) bool {
	return int(ou.Role) >= int(targetRole) && ou.Role >= OrganizationUserRoleAdmin
}
