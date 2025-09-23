package organization

import (
	"testing"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestOrganization_CanChangeRole(t *testing.T) {
	orgID := shared.NewUUID[Organization]()
	userID := shared.NewUUID[user.User]()
	org := NewOrganization(orgID, OrganizationTypeNormal, "Test Org", "TEST", lo.ToPtr("https://kotohiro.com/aaaa"), userID)

	tests := []struct {
		name            string
		currentUserRole OrganizationUserRole
		targetRole      OrganizationUserRole
		expected        bool
	}{
		{
			name:            "SuperAdmin can change to any role",
			currentUserRole: OrganizationUserRoleSuperAdmin,
			targetRole:      OrganizationUserRoleMember,
			expected:        true,
		},
		{
			name:            "SuperAdmin can change to SuperAdmin",
			currentUserRole: OrganizationUserRoleSuperAdmin,
			targetRole:      OrganizationUserRoleSuperAdmin,
			expected:        true,
		},
		{
			name:            "Owner can change to Member",
			currentUserRole: OrganizationUserRoleOwner,
			targetRole:      OrganizationUserRoleMember,
			expected:        true,
		},
		{
			name:            "Owner can change to Admin",
			currentUserRole: OrganizationUserRoleOwner,
			targetRole:      OrganizationUserRoleAdmin,
			expected:        true,
		},
		{
			name:            "Owner can change to Owner",
			currentUserRole: OrganizationUserRoleOwner,
			targetRole:      OrganizationUserRoleOwner,
			expected:        true,
		},
		{
			name:            "Owner cannot change to SuperAdmin",
			currentUserRole: OrganizationUserRoleOwner,
			targetRole:      OrganizationUserRoleSuperAdmin,
			expected:        false,
		},
		{
			name:            "Admin can change to Member",
			currentUserRole: OrganizationUserRoleAdmin,
			targetRole:      OrganizationUserRoleMember,
			expected:        true,
		},
		{
			name:            "Admin can change to Admin",
			currentUserRole: OrganizationUserRoleAdmin,
			targetRole:      OrganizationUserRoleAdmin,
			expected:        true,
		},
		{
			name:            "Admin cannot change to Owner",
			currentUserRole: OrganizationUserRoleAdmin,
			targetRole:      OrganizationUserRoleOwner,
			expected:        false,
		},
		{
			name:            "Admin cannot change to SuperAdmin",
			currentUserRole: OrganizationUserRoleAdmin,
			targetRole:      OrganizationUserRoleSuperAdmin,
			expected:        false,
		},
		{
			name:            "Member cannot change any role",
			currentUserRole: OrganizationUserRoleMember,
			targetRole:      OrganizationUserRoleMember,
			expected:        false,
		},
		{
			name:            "Member cannot change to Admin",
			currentUserRole: OrganizationUserRoleMember,
			targetRole:      OrganizationUserRoleAdmin,
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := org.CanChangeRole(tt.currentUserRole, tt.targetRole)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewOrganization(t *testing.T) {
	orgID := shared.NewUUID[Organization]()
	userID := shared.NewUUID[user.User]()

	tests := []struct {
		name    string
		orgType OrganizationType
		orgName string
		code    string
		iconURL *string
		ownerID shared.UUID[user.User]
	}{
		{
			name:    "Normal organization",
			orgType: OrganizationTypeNormal,
			orgName: "Test Organization",
			code:    "TEST123",
			iconURL: lo.ToPtr("https://kotohiro.com/aaaa"),
			ownerID: userID,
		},
		{
			name:    "Government organization",
			orgType: OrganizationTypeGovernment,
			orgName: "Government Agency",
			code:    "GOV001",
			iconURL: lo.ToPtr("https://kotohiro.com/aaaa"),
			ownerID: userID,
		},
		{
			name:    "Councillor organization",
			orgType: OrganizationTypeCouncillor,
			orgName: "City Council",
			code:    "COUNCIL1",
			iconURL: lo.ToPtr("https://kotohiro.com/aaaa"),
			ownerID: userID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org := NewOrganization(orgID, tt.orgType, tt.orgName, tt.code, tt.iconURL, tt.ownerID)

			assert.Equal(t, orgID, org.OrganizationID)
			assert.Equal(t, tt.orgType, org.OrganizationType)
			assert.Equal(t, tt.orgName, org.Name)
			assert.Equal(t, tt.code, org.Code)
			assert.Equal(t, tt.iconURL, org.IconURL)
			assert.Equal(t, tt.ownerID, org.OwnerID)
		})
	}
}
