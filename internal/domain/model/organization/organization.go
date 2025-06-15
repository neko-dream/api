package organization

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type OrganizationRepository interface {
	// 組織の取得
	FindByID(ctx context.Context, id shared.UUID[Organization]) (*Organization, error)
	FindByIDs(ctx context.Context, ids []shared.UUID[Organization]) ([]*Organization, error)
	FindByName(ctx context.Context, name string) (*Organization, error)
	FindByCode(ctx context.Context, code string) (*Organization, error)

	// 組織の作成・更新・削除
	Create(ctx context.Context, org *Organization) error
}

type OrganizationType int

const (
	OrganizationTypeNormal OrganizationType = iota + 1
	OrganizationTypeGovernment
	OrganizationTypeCouncillor
)

type Organization struct {
	OrganizationID   shared.UUID[Organization]
	OrganizationType OrganizationType
	Name             string
	Code             string
	OwnerID          shared.UUID[user.User]
}

func NewOrganization(
	organizationID shared.UUID[Organization],
	organizationType OrganizationType,
	name string,
	code string,
	ownerID shared.UUID[user.User],
) *Organization {
	return &Organization{
		OrganizationID:   organizationID,
		OrganizationType: organizationType,
		Name:             name,
		Code:             code,
		OwnerID:          ownerID,
	}
}

// CanChangeRole はユーザーが他のユーザーのロールを変更できるかを判断するのじゃ
func (o *Organization) CanChangeRole(currentUserRole OrganizationUserRole, targetRole OrganizationUserRole) bool {
	// 上位のロールにのみロールを割り当てることができるのじゃ
	return currentUserRole >= OrganizationUserRoleAdmin && int(currentUserRole) >= int(targetRole)
}
