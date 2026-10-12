package organization

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/neko-dream/api/internal/domain/messages"
	model "github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

type loginMembers struct {
	model.OrganizationUserRepository
	rows []*model.OrganizationUser
	err  error
}

func (r *loginMembers) FindByUserID(context.Context, shared.UUID[user.User]) ([]*model.OrganizationUser, error) {
	return r.rows, r.err
}

type loginOrgs struct {
	model.OrganizationRepository
	org *model.Organization
	err error
}

func (r *loginOrgs) FindByID(context.Context, shared.UUID[model.Organization]) (*model.Organization, error) {
	return r.org, r.err
}

func TestResolveLoginOrganization(t *testing.T) {
	uid := shared.NewUUID[user.User]()
	oid := shared.NewUUID[model.Organization]()
	requested := shared.UUID[any](oid.UUID())
	wrong := shared.NewUUID[any]()
	membership := &model.OrganizationUser{OrganizationID: oid, UserID: uid, Role: model.OrganizationUserRoleOwner}
	dbErr := errors.New("membership database unavailable")
	for _, tc := range []struct {
		name                       string
		rows                       []*model.OrganizationUser
		requested                  *shared.UUID[any]
		memberErr, orgErr, wantErr error
		wantOrg                    bool
	}{
		{"ordinary account", nil, nil, nil, nil, nil, false},
		{"automatic organization", []*model.OrganizationUser{membership}, nil, nil, nil, nil, true},
		{"same requested organization", []*model.OrganizationUser{membership}, &requested, nil, nil, nil, true},
		{"other requested organization", []*model.OrganizationUser{membership}, &wrong, nil, nil, messages.OrganizationPermissionDenied, false},
		{"code alone not membership", nil, &requested, nil, nil, messages.OrganizationPermissionDenied, false},
		{"multiple memberships fail closed", []*model.OrganizationUser{membership, membership}, nil, nil, nil, messages.OrganizationMembershipConflict, false},
		{"nil membership", []*model.OrganizationUser{nil}, nil, nil, nil, messages.OrganizationPermissionDenied, false},
		{"membership lookup error", nil, nil, dbErr, nil, dbErr, false},
		{"missing organization", []*model.OrganizationUser{membership}, nil, nil, sql.ErrNoRows, messages.OrganizationPermissionDenied, false},
		{"organization lookup error", []*model.OrganizationUser{membership}, nil, nil, dbErr, dbErr, false},
		{"invalid role", []*model.OrganizationUser{{OrganizationID: oid, UserID: uid, Role: 0}}, nil, nil, nil, messages.OrganizationPermissionDenied, false},
		{"other user membership", []*model.OrganizationUser{{OrganizationID: oid, UserID: shared.NewUUID[user.User](), Role: model.OrganizationUserRoleOwner}}, nil, nil, nil, messages.OrganizationPermissionDenied, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ResolveLoginOrganization(context.Background(), &loginMembers{rows: tc.rows, err: tc.memberErr}, &loginOrgs{org: &model.Organization{OrganizationID: oid}, err: tc.orgErr}, uid, tc.requested)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error %v want %v", err, tc.wantErr)
			}
			if (got != nil) != tc.wantOrg || (got != nil && *got != requested) {
				t.Fatalf("wrong organization %v", got)
			}
		})
	}
}
