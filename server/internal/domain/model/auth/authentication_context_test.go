package auth

import (
	"testing"

	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationContext_IsInOrganization(t *testing.T) {
	tests := []struct {
		name           string
		organizationID *shared.UUID[organization.Organization]
		expected       bool
	}{
		{
			name:           "組織IDがある場合はtrue",
			organizationID: lo.ToPtr(organization.KotohiroOrganizationID),
			expected:       true,
		},
		{
			name:           "組織IDがnilの場合はfalse",
			organizationID: nil,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &AuthenticationContext{
				OrganizationID: tt.organizationID,
			}
			assert.Equal(t, tt.expected, ac.IsInOrganization())
		})
	}
}

func TestAuthenticationContext_HasOrganizationRole(t *testing.T) {
	tests := []struct {
		name             string
		organizationID   *shared.UUID[organization.Organization]
		organizationRole *organization.OrganizationUserRole
		minRole          organization.OrganizationUserRole
		expected         bool
	}{
		{
			name:             "SuperAdminはSuperAdmin以上の権限を持つ",
			organizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
			organizationRole: lo.ToPtr(organization.OrganizationUserRoleSuperAdmin),
			minRole:          organization.OrganizationUserRoleSuperAdmin,
			expected:         true,
		},
		{
			name:             "SuperAdminはMember以上の権限を持つ",
			organizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
			organizationRole: lo.ToPtr(organization.OrganizationUserRoleSuperAdmin),
			minRole:          organization.OrganizationUserRoleMember,
			expected:         true,
		},
		{
			name:             "OwnerはOwner以上の権限を持つ",
			organizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
			organizationRole: lo.ToPtr(organization.OrganizationUserRoleOwner),
			minRole:          organization.OrganizationUserRoleOwner,
			expected:         true,
		},
		{
			name:             "OwnerはSuperAdmin以上の権限を持たない",
			organizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
			organizationRole: lo.ToPtr(organization.OrganizationUserRoleOwner),
			minRole:          organization.OrganizationUserRoleSuperAdmin,
			expected:         false,
		},
		{
			name:             "AdminはAdmin以上の権限を持つ",
			organizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
			organizationRole: lo.ToPtr(organization.OrganizationUserRoleAdmin),
			minRole:          organization.OrganizationUserRoleAdmin,
			expected:         true,
		},
		{
			name:             "AdminはOwner以上の権限を持たない",
			organizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
			organizationRole: lo.ToPtr(organization.OrganizationUserRoleAdmin),
			minRole:          organization.OrganizationUserRoleOwner,
			expected:         false,
		},
		{
			name:             "MemberはMember以上の権限を持つ",
			organizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
			organizationRole: lo.ToPtr(organization.OrganizationUserRoleMember),
			minRole:          organization.OrganizationUserRoleMember,
			expected:         true,
		},
		{
			name:             "MemberはAdmin以上の権限を持たない",
			organizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
			organizationRole: lo.ToPtr(organization.OrganizationUserRoleMember),
			minRole:          organization.OrganizationUserRoleAdmin,
			expected:         false,
		},
		{
			name:             "組織に所属していない場合はfalse",
			organizationID:   nil,
			organizationRole: lo.ToPtr(organization.OrganizationUserRoleSuperAdmin),
			minRole:          organization.OrganizationUserRoleMember,
			expected:         false,
		},
		{
			name:             "ロールがnilの場合はfalse",
			organizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
			organizationRole: nil,
			minRole:          organization.OrganizationUserRoleMember,
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &AuthenticationContext{
				OrganizationID:   tt.organizationID,
				OrganizationRole: tt.organizationRole,
			}
			assert.Equal(t, tt.expected, ac.HasOrganizationRole(tt.minRole))
		})
	}
}

func TestAuthenticationContext_RoleChecks(t *testing.T) {
	t.Run("IsSuperAdmin", func(t *testing.T) {
		tests := []struct {
			name     string
			role     *organization.OrganizationUserRole
			expected bool
		}{
			{
				name:     "SuperAdminはtrue",
				role:     lo.ToPtr(organization.OrganizationUserRoleSuperAdmin),
				expected: true,
			},
			{
				name:     "Ownerはfalse",
				role:     lo.ToPtr(organization.OrganizationUserRoleOwner),
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ac := &AuthenticationContext{
					OrganizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
					OrganizationRole: tt.role,
				}
				assert.Equal(t, tt.expected, ac.IsSuperAdmin())
			})
		}
	})

	t.Run("IsOwner", func(t *testing.T) {
		tests := []struct {
			name     string
			role     *organization.OrganizationUserRole
			expected bool
		}{
			{
				name:     "SuperAdminはtrue",
				role:     lo.ToPtr(organization.OrganizationUserRoleSuperAdmin),
				expected: true,
			},
			{
				name:     "Ownerはtrue",
				role:     lo.ToPtr(organization.OrganizationUserRoleOwner),
				expected: true,
			},
			{
				name:     "Adminはfalse",
				role:     lo.ToPtr(organization.OrganizationUserRoleAdmin),
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ac := &AuthenticationContext{
					OrganizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
					OrganizationRole: tt.role,
				}
				assert.Equal(t, tt.expected, ac.IsOwner())
			})
		}
	})

	t.Run("IsAdmin", func(t *testing.T) {
		tests := []struct {
			name     string
			role     *organization.OrganizationUserRole
			expected bool
		}{
			{
				name:     "SuperAdminはtrue",
				role:     lo.ToPtr(organization.OrganizationUserRoleSuperAdmin),
				expected: true,
			},
			{
				name:     "Ownerはtrue",
				role:     lo.ToPtr(organization.OrganizationUserRoleOwner),
				expected: true,
			},
			{
				name:     "Adminはtrue",
				role:     lo.ToPtr(organization.OrganizationUserRoleAdmin),
				expected: true,
			},
			{
				name:     "Memberはfalse",
				role:     lo.ToPtr(organization.OrganizationUserRoleMember),
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ac := &AuthenticationContext{
					OrganizationID:   lo.ToPtr(organization.KotohiroOrganizationID),
					OrganizationRole: tt.role,
				}
				assert.Equal(t, tt.expected, ac.IsAdmin())
			})
		}
	})
}

func TestAuthenticationContext_RequiresPasswordChange(t *testing.T) {
	tests := []struct {
		name                   string
		requiredPasswordChange bool
		expected               bool
	}{
		{
			name:                   "パスワード変更が必要な場合はtrue",
			requiredPasswordChange: true,
			expected:               true,
		},
		{
			name:                   "パスワード変更が不要な場合はfalse",
			requiredPasswordChange: false,
			expected:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &AuthenticationContext{
				RequiredPasswordChange: tt.requiredPasswordChange,
			}
			assert.Equal(t, tt.expected, ac.RequiresPasswordChange())
		})
	}
}

func TestAuthenticationContext_IsKotohiro(t *testing.T) {
	otherOrgID := shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000002")

	tests := []struct {
		name           string
		organizationID *shared.UUID[organization.Organization]
		expected       bool
	}{
		{
			name:           "Kotohiro組織の場合はtrue",
			organizationID: lo.ToPtr(organization.KotohiroOrganizationID),
			expected:       true,
		},
		{
			name:           "他の組織の場合はfalse",
			organizationID: &otherOrgID,
			expected:       false,
		},
		{
			name:           "組織IDがnilの場合はfalse",
			organizationID: nil,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &AuthenticationContext{
				OrganizationID: tt.organizationID,
			}
			assert.Equal(t, tt.expected, ac.IsKotohiro())
		})
	}
}
