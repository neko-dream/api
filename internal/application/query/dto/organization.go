package dto

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
)

type Organization struct {
	OrganizationID   string    `json:"organization_id"`   // 組織ID
	Name             string    `json:"name"`              // 組織名
	Code             string    `json:"code"`              // 組織コード
	IconURL          *string   `json:"icon_url"`          // 組織アイコンURL
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
	User             User             `json:"user"`
}

func (o *OrganizationResponse) ToResponse() oas.Organization {
	return oas.Organization{
		ID:       o.Organization.OrganizationID,
		Name:     o.Organization.Name,
		Code:     o.Organization.Code,
		Type:     o.Organization.OrganizationType,
		Role:     o.OrganizationUser.Role,
		RoleName: o.OrganizationUser.RoleName,
		IconURL:  utils.ToOptNil[oas.OptNilString](o.Organization.IconURL),
	}
}

func (o *OrganizationResponse) ToUserResponse() oas.OrganizationUser {
	user := o.User.ToResponse()
	return oas.OrganizationUser{
		DisplayName: user.DisplayName,
		DisplayID:   user.DisplayID,
		IconURL:     user.IconURL,
		Role:        o.OrganizationUser.Role,
		RoleName:    o.OrganizationUser.RoleName,
		UserID:      o.OrganizationUser.UserID,
	}
}

type OrganizationAlias struct {
	AliasID   string     `json:"alias_id"`   // エイリアスID
	AliasName string     `json:"alias_name"` // エイリアス名
	CreatedAt *time.Time `json:"created_at"` // エイリアスの作成日時
}

func (o *OrganizationAlias) ToResponse() oas.OrganizationAlias {
	if o == nil {
		return oas.OrganizationAlias{}
	}
	if o.AliasName == "" {
		o.AliasName = ""
	}

	var createdAt oas.OptNilString
	if o.CreatedAt != nil {
		createdAt = utils.ToOptNil[oas.OptNilString](o.CreatedAt.Format(time.RFC3339))
	}
	return oas.OrganizationAlias{
		AliasID:   o.AliasID,
		AliasName: o.AliasName,
		CreatedAt: createdAt,
	}
}
