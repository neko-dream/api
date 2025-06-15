package organization

import (
	"testing"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/stretchr/testify/assert"
)

func TestNewOrganizationUserRole(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected OrganizationUserRole
	}{
		{
			name:     "Valid Member role",
			input:    int(OrganizationUserRoleMember),
			expected: OrganizationUserRoleMember,
		},
		{
			name:     "Valid Admin role",
			input:    int(OrganizationUserRoleAdmin),
			expected: OrganizationUserRoleAdmin,
		},
		{
			name:     "Valid Owner role",
			input:    int(OrganizationUserRoleOwner),
			expected: OrganizationUserRoleOwner,
		},
		{
			name:     "Valid SuperAdmin role",
			input:    int(OrganizationUserRoleSuperAdmin),
			expected: OrganizationUserRoleSuperAdmin,
		},
		{
			name:     "Invalid role below range defaults to Member",
			input:    0,
			expected: OrganizationUserRoleMember,
		},
		{
			name:     "Invalid role above range defaults to Member",
			input:    100,
			expected: OrganizationUserRoleMember,
		},
		{
			name:     "Negative role defaults to Member",
			input:    -1,
			expected: OrganizationUserRoleMember,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewOrganizationUserRole(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoleToName(t *testing.T) {
	tests := []struct {
		name     string
		role     OrganizationUserRole
		expected string
	}{
		{
			name:     "Member role",
			role:     OrganizationUserRoleMember,
			expected: "メンバー",
		},
		{
			name:     "Admin role",
			role:     OrganizationUserRoleAdmin,
			expected: "管理者",
		},
		{
			name:     "Owner role",
			role:     OrganizationUserRoleOwner,
			expected: "オーナー",
		},
		{
			name:     "SuperAdmin role",
			role:     OrganizationUserRoleSuperAdmin,
			expected: "運営",
		},
		{
			name:     "Invalid role defaults to Member",
			role:     OrganizationUserRole(100),
			expected: "メンバー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoleToName(tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOrganizationUser_SetRole(t *testing.T) {
	orgUserID := shared.NewUUID[OrganizationUser]()
	orgID := shared.NewUUID[Organization]()
	userID := shared.NewUUID[user.User]()

	tests := []struct {
		name        string
		initialRole OrganizationUserRole
		newRole     OrganizationUserRole
		expectError bool
	}{
		{
			name:        "Set valid Member role",
			initialRole: OrganizationUserRoleAdmin,
			newRole:     OrganizationUserRoleMember,
			expectError: false,
		},
		{
			name:        "Set valid Admin role",
			initialRole: OrganizationUserRoleMember,
			newRole:     OrganizationUserRoleAdmin,
			expectError: false,
		},
		{
			name:        "Set valid Owner role",
			initialRole: OrganizationUserRoleAdmin,
			newRole:     OrganizationUserRoleOwner,
			expectError: false,
		},
		{
			name:        "Set valid SuperAdmin role",
			initialRole: OrganizationUserRoleOwner,
			newRole:     OrganizationUserRoleSuperAdmin,
			expectError: false,
		},
		{
			name:        "Set invalid role below range",
			initialRole: OrganizationUserRoleMember,
			newRole:     OrganizationUserRole(0),
			expectError: true,
		},
		{
			name:        "Set invalid role above range",
			initialRole: OrganizationUserRoleMember,
			newRole:     OrganizationUserRole(100),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgUser := NewOrganizationUser(orgUserID, orgID, userID, tt.initialRole)
			err := orgUser.SetRole(tt.newRole)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.initialRole, orgUser.Role) // Role should not change on error
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newRole, orgUser.Role) // Role should be updated
			}
		})
	}
}

func TestOrganizationUser_HasPermissionToChangeRoleTo(t *testing.T) {
	orgUserID := shared.NewUUID[OrganizationUser]()
	orgID := shared.NewUUID[Organization]()
	userID := shared.NewUUID[user.User]()

	tests := []struct {
		name       string
		userRole   OrganizationUserRole
		targetRole OrganizationUserRole
		expected   bool
	}{
		// SuperAdmin tests
		{
			name:       "SuperAdmin can change to Member",
			userRole:   OrganizationUserRoleSuperAdmin,
			targetRole: OrganizationUserRoleMember,
			expected:   true,
		},
		{
			name:       "SuperAdmin can change to Admin",
			userRole:   OrganizationUserRoleSuperAdmin,
			targetRole: OrganizationUserRoleAdmin,
			expected:   true,
		},
		{
			name:       "SuperAdmin can change to Owner",
			userRole:   OrganizationUserRoleSuperAdmin,
			targetRole: OrganizationUserRoleOwner,
			expected:   true,
		},
		{
			name:       "SuperAdmin can change to SuperAdmin",
			userRole:   OrganizationUserRoleSuperAdmin,
			targetRole: OrganizationUserRoleSuperAdmin,
			expected:   true,
		},
		// Owner tests
		{
			name:       "Owner can change to Member",
			userRole:   OrganizationUserRoleOwner,
			targetRole: OrganizationUserRoleMember,
			expected:   true,
		},
		{
			name:       "Owner can change to Admin",
			userRole:   OrganizationUserRoleOwner,
			targetRole: OrganizationUserRoleAdmin,
			expected:   true,
		},
		{
			name:       "Owner can change to Owner",
			userRole:   OrganizationUserRoleOwner,
			targetRole: OrganizationUserRoleOwner,
			expected:   true,
		},
		{
			name:       "Owner cannot change to SuperAdmin",
			userRole:   OrganizationUserRoleOwner,
			targetRole: OrganizationUserRoleSuperAdmin,
			expected:   false,
		},
		// Admin tests
		{
			name:       "Admin can change to Member",
			userRole:   OrganizationUserRoleAdmin,
			targetRole: OrganizationUserRoleMember,
			expected:   true,
		},
		{
			name:       "Admin can change to Admin",
			userRole:   OrganizationUserRoleAdmin,
			targetRole: OrganizationUserRoleAdmin,
			expected:   true,
		},
		{
			name:       "Admin cannot change to Owner",
			userRole:   OrganizationUserRoleAdmin,
			targetRole: OrganizationUserRoleOwner,
			expected:   false,
		},
		{
			name:       "Admin cannot change to SuperAdmin",
			userRole:   OrganizationUserRoleAdmin,
			targetRole: OrganizationUserRoleSuperAdmin,
			expected:   false,
		},
		// Member tests
		{
			name:       "Member cannot change to Member",
			userRole:   OrganizationUserRoleMember,
			targetRole: OrganizationUserRoleMember,
			expected:   false,
		},
		{
			name:       "Member cannot change to Admin",
			userRole:   OrganizationUserRoleMember,
			targetRole: OrganizationUserRoleAdmin,
			expected:   false,
		},
		{
			name:       "Member cannot change to Owner",
			userRole:   OrganizationUserRoleMember,
			targetRole: OrganizationUserRoleOwner,
			expected:   false,
		},
		{
			name:       "Member cannot change to SuperAdmin",
			userRole:   OrganizationUserRoleMember,
			targetRole: OrganizationUserRoleSuperAdmin,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgUser := NewOrganizationUser(orgUserID, orgID, userID, tt.userRole)
			result := orgUser.HasPermissionToChangeRoleTo(tt.targetRole)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewOrganizationUser(t *testing.T) {
	orgUserID := shared.NewUUID[OrganizationUser]()
	orgID := shared.NewUUID[Organization]()
	userID := shared.NewUUID[user.User]()

	tests := []struct {
		name string
		role OrganizationUserRole
	}{
		{
			name: "Create Member",
			role: OrganizationUserRoleMember,
		},
		{
			name: "Create Admin",
			role: OrganizationUserRoleAdmin,
		},
		{
			name: "Create Owner",
			role: OrganizationUserRoleOwner,
		},
		{
			name: "Create SuperAdmin",
			role: OrganizationUserRoleSuperAdmin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgUser := NewOrganizationUser(orgUserID, orgID, userID, tt.role)

			assert.Equal(t, orgUserID, orgUser.OrganizationUserID)
			assert.Equal(t, orgID, orgUser.OrganizationID)
			assert.Equal(t, userID, orgUser.UserID)
			assert.Equal(t, tt.role, orgUser.Role)
		})
	}
}
