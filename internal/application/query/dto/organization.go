package dto

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/organization"
)

type Organization struct {
	ID               string    `json:"id"`                // 組織ID
	Name             string    `json:"name"`              // 組織名
	Code             string    `json:"code"`              // 組織コード
	OrganizationType int       `json:"organization_type"` // 組織の種類
	CreatedAt        time.Time `json:"created_at"`        // 組織の作成日時
	UpdatedAt        time.Time `json:"updated_at"`        // 組織の更新日時
}

type OrganizationUser struct {
	OrganizationUserID string `json:"organization_user_id"` // 組織ユーザーID
	OrganizationID     string `json:"organization_id"`      // 組織ID
	UserID             string `json:"user_id"`              // ユーザーID
	Role               int    `json:"role"`                 // 組織のユーザーのロール
	RoleName           string `json:"role_name"`            // 組織のユーザーのロール名
}

func (ou *OrganizationUser) SetRoleName(role int) {
	ou.RoleName = organization.RoleToName(organization.OrganizationUserRole(role))
}

type OrganizationResponse struct {
	Organization     Organization     `json:"organization"`
	OrganizationUser OrganizationUser `json:"organization_user"`
}
