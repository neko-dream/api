package organization

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/api/internal/domain/messages"
	model "github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

// ResolveLoginOrganization trusts persisted membership, never a login-page or
// organization code alone. Unaffiliated users keep their ordinary session.
func ResolveLoginOrganization(ctx context.Context, members model.OrganizationUserRepository, orgs model.OrganizationRepository, userID shared.UUID[user.User], requested *shared.UUID[any]) (*shared.UUID[any], error) {
	memberships, err := members.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(memberships) > 1 {
		return nil, messages.OrganizationMembershipConflict
	}
	if len(memberships) == 0 {
		if requested != nil && !requested.IsZero() {
			return nil, messages.OrganizationPermissionDenied
		}
		return nil, nil
	}
	membership := memberships[0]
	if membership == nil || membership.UserID != userID || membership.OrganizationID.IsZero() {
		return nil, messages.OrganizationPermissionDenied
	}
	if !membership.Role.IsValid() {
		return nil, messages.OrganizationPermissionDenied
	}
	id := shared.UUID[any](membership.OrganizationID.UUID())
	if requested != nil && !requested.IsZero() && *requested != id {
		return nil, messages.OrganizationPermissionDenied
	}
	org, err := orgs.FindByID(ctx, membership.OrganizationID)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && org == nil) {
		return nil, messages.OrganizationPermissionDenied
	}
	if err != nil {
		return nil, err
	}
	if org.OrganizationID != membership.OrganizationID {
		return nil, messages.OrganizationPermissionDenied
	}
	return &id, nil
}
