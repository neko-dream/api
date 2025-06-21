package auth

import (
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type AuthenticationContext struct {
	UserID    shared.UUID[user.User]
	SessionID shared.UUID[session.Session]

	// ユーザープロフィール情報
	DisplayName     *string
	DisplayID       *string
	IconURL         *string
	IsRegistered    bool
	IsEmailVerified bool

	// パスワード変更要求
	RequiredPasswordChange bool

	OrganizationID   *string
	OrganizationCode *string
	OrganizationRole *shared.OrganizationUserRole
}

// 組織コンテキスト内かどうかを確認
func (ac *AuthenticationContext) IsInOrganization() bool {
	return ac.OrganizationID != nil
}

// 指定された役割以上の権限を持つかを確認
func (ac *AuthenticationContext) HasOrganizationRole(minRole shared.OrganizationUserRole) bool {
	if !ac.IsInOrganization() || ac.OrganizationRole == nil {
		return false
	}
	// 数値が小さいほど高い権限 (SuperAdmin=10, Owner=20, Admin=30, Member=40)
	return int(*ac.OrganizationRole) <= int(minRole)
}

// スーパー管理者権限を持つかを確認
func (ac *AuthenticationContext) IsSuperAdmin() bool {
	return ac.HasOrganizationRole(shared.OrganizationUserRoleSuperAdmin)
}

// オーナー権限を持つかを確認
func (ac *AuthenticationContext) IsOwner() bool {
	return ac.HasOrganizationRole(shared.OrganizationUserRoleOwner)
}

// 管理者以上の権限を持つかを確認
func (ac *AuthenticationContext) IsAdmin() bool {
	return ac.HasOrganizationRole(shared.OrganizationUserRoleAdmin)
}

// パスワード変更が必要かを確認
func (ac *AuthenticationContext) RequiresPasswordChange() bool {
	return ac.RequiredPasswordChange
}
