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

func NewOrganizationUserRole(role int) OrganizationUserRole {
	if role < int(OrganizationUserRoleMember) || role > int(OrganizationUserRoleSuperAdmin) {
		return OrganizationUserRoleMember
	}
	if role == 0 {
		return OrganizationUserRoleMember
	}

	return OrganizationUserRole(role)
}

func RoleToName(role OrganizationUserRole) string {
	switch role {
	case OrganizationUserRoleMember:
		return "メンバー"
	case OrganizationUserRoleAdmin:
		return "管理者"
	case OrganizationUserRoleOwner:
		return "オーナー"
	case OrganizationUserRoleSuperAdmin:
		return "運営"
	default:
		return "メンバー"
	}
}

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

// SetRole
func (ou *OrganizationUser) SetRole(role OrganizationUserRole) error {
	if role < OrganizationUserRoleMember || role > OrganizationUserRoleSuperAdmin {
		return errors.New("invalid role")
	}
	ou.Role = role
	return nil
}

// HasPermissionToChangeRoleTo
func (ou *OrganizationUser) HasPermissionToChangeRoleTo(targetRole OrganizationUserRole) bool {
	return int(ou.Role) >= int(targetRole) && ou.Role >= OrganizationUserRoleAdmin
}
