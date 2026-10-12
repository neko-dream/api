package talksession_usecase

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
)

type creationUserStub struct {
	user.UserRepository
	value *user.User
	err   error
}

func (s creationUserStub) FindByID(context.Context, shared.UUID[user.User]) (*user.User, error) {
	return s.value, s.err
}

type creationGrantStub struct {
	granted bool
	err     error
}

func (s *creationGrantStub) Get(context.Context, shared.UUID[user.User]) (bool, error) {
	return s.granted, s.err
}
func (s *creationGrantStub) Set(context.Context, shared.UUID[user.User], bool) error {
	panic("creation must not grant permissions")
}

type creationOrganizationStub struct {
	organization.OrganizationUserRepository
	calls int
}

func (s *creationOrganizationStub) FindByOrganizationIDAndUserID(context.Context, shared.UUID[organization.Organization], shared.UUID[user.User]) (*organization.OrganizationUser, error) {
	s.calls++
	return nil, nil
}

type creationSessionStub struct {
	talksession.TalkSessionRepository
	created *talksession.TalkSession
}

func (s *creationSessionStub) Create(_ context.Context, value *talksession.TalkSession) error {
	s.created = value
	return nil
}

func creationUser(provider shared.AuthProviderName, withdrawn bool) user.User {
	name := "creator"
	id := shared.NewUUID[user.User]()
	return user.NewUserWithWithdrawalDate(id, &name, &name, id.String(), provider, nil, sql.NullTime{Time: time.Now(), Valid: withdrawn})
}

func TestCreationRequiresIndependentGrant(t *testing.T) {
	for _, tc := range []struct {
		name                      string
		granted, guest, withdrawn bool
		want                      error
	}{
		{"ungranted", false, false, false, messages.ForbiddenError},
		{"granted without organization", true, false, false, messages.TalkSessionThemeTooLong},
		{"guest even if granted", true, true, false, messages.GuestCannotCreateSessionError},
		{"withdrawn even if granted", true, false, true, messages.ForbiddenError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			provider := shared.ProviderGoogle
			if tc.guest {
				provider = shared.ProviderGuest
			}
			u := creationUser(provider, tc.withdrawn)
			orgs := &creationOrganizationStub{}
			h := &startTalkSessionHandler{UserRepository: creationUserStub{value: &u}, creationPermissions: &creationGrantStub{granted: tc.granted}, OrganizationUserRepository: orgs}
			_, err := h.Execute(context.Background(), StartTalkSessionUseCaseInput{OwnerID: u.UserID(), Theme: strings.Repeat("x", 101)})
			if !errors.Is(err, tc.want) || orgs.calls != 0 {
				t.Fatalf("err=%v org calls=%d", err, orgs.calls)
			}
		})
	}
	u := creationUser(shared.ProviderGoogle, false)
	failure := errors.New("permission database unavailable")
	h := &startTalkSessionHandler{UserRepository: creationUserStub{value: &u}, creationPermissions: &creationGrantStub{err: failure}}
	if _, err := h.Execute(context.Background(), StartTalkSessionUseCaseInput{OwnerID: u.UserID()}); !errors.Is(err, failure) {
		t.Fatal("permission read failure did not fail closed", err)
	}
	h.UserRepository = creationUserStub{}
	if _, err := h.Execute(context.Background(), StartTalkSessionUseCaseInput{}); !errors.Is(err, messages.UserNotFoundError) {
		t.Fatal("missing user accepted", err)
	}
}

func TestGrantedUserCreatesWithoutOrganizationAndRevocationBlocksNextRequest(t *testing.T) {
	url := os.Getenv("ADMIN_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("ADMIN_TEST_DATABASE_URL required")
	}
	c, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	u := creationUser(shared.ProviderGoogle, false)
	grants := &creationGrantStub{granted: true}
	sessions := &creationSessionStub{}
	orgs := &creationOrganizationStub{}
	h := &startTalkSessionHandler{UserRepository: creationUserStub{value: &u}, creationPermissions: grants, TalkSessionRepository: sessions, OrganizationUserRepository: orgs, DBManager: db.NewDBManager(c)}
	input := StartTalkSessionUseCaseInput{OwnerID: u.UserID(), Theme: "permission test", ScheduledEndTime: time.Now().Add(time.Hour)}
	if _, err := h.Execute(context.Background(), input); err != nil || sessions.created == nil || sessions.created.OrganizationID() != nil || orgs.calls != 0 {
		t.Fatalf("creation err=%v org calls=%d", err, orgs.calls)
	}
	grants.granted = false
	sessions.created = nil
	if _, err := h.Execute(context.Background(), input); !errors.Is(err, messages.ForbiddenError) || sessions.created != nil {
		t.Fatal("revocation did not block next request", err)
	}
	grants.granted = true
	orgID := shared.NewUUID[organization.Organization]().String()
	input.SessionClaim = &session.Claim{OrganizationID: &orgID}
	if _, err := h.Execute(context.Background(), input); !errors.Is(err, messages.ForbiddenError) || sessions.created != nil || orgs.calls != 1 {
		t.Fatal("unrelated organization identity accepted", err)
	}
}
