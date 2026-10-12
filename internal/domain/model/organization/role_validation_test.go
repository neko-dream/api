package organization

import "testing"

func TestOrganizationUserRoleValidation(t *testing.T) {
	for _, role := range []OrganizationUserRole{-1, 0, 9, 10, 11, 19, 20, 21, 29, 30, 31, 39, 40, 41, 100} {
		want := role == 10 || role == 20 || role == 30 || role == 40
		if role.IsValid() != want {
			t.Errorf("role %d: validity must be %v", role, want)
		}
		got := NewOrganizationUserRole(int(role))
		if want && got != role || !want && got != OrganizationUserRoleMember {
			t.Errorf("role %d: unexpected constructor result %d", role, got)
		}
		member := OrganizationUser{Role: OrganizationUserRoleOwner}
		if err := member.SetRole(role); (err == nil) != want {
			t.Errorf("role %d: SetRole error %v", role, err)
		}
		if !want && member.Role != OrganizationUserRoleOwner {
			t.Errorf("invalid role %d mutated member", role)
		}
		actor := OrganizationUser{Role: role}
		if !want && actor.HasPermissionToChangeRoleTo(OrganizationUserRoleMember) {
			t.Errorf("invalid actor role %d granted permission", role)
		}
		owner := OrganizationUser{Role: OrganizationUserRoleOwner}
		if !want && owner.HasPermissionToChangeRoleTo(role) {
			t.Errorf("invalid target role %d accepted", role)
		}
	}
}
