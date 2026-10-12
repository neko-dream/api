package handler

import (
	"context"
	"errors"
	"testing"

	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/auth"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/service"
)

type currentAuthStub struct {
	service.AuthorizationService
	id  shared.UUID[user.User]
	err error
}

func (s currentAuthStub) RequireAuthentication(context.Context) (*auth.AuthenticationContext, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &auth.AuthenticationContext{UserID: s.id}, nil
}

type currentQueryStub struct {
	out   *userquery.CurrentOutput
	id    shared.UUID[user.User]
	calls int
	err   error
}

func (s *currentQueryStub) FindByID(_ context.Context, id shared.UUID[user.User]) (*userquery.CurrentOutput, error) {
	s.id = id
	s.calls++
	return s.out, s.err
}

func TestCurrentUserReadsSelfAndCurrentPermission(t *testing.T) {
	id := shared.NewUUID[user.User]()
	q := &currentQueryStub{out: &userquery.CurrentOutput{CanCreateTalkSession: true}}
	h := &userHandler{authorizationService: currentAuthStub{id: id}, currentUser: q}
	for _, allowed := range []bool{true, false} {
		q.out.CanCreateTalkSession = allowed
		out, err := h.GetCurrentUser(context.Background())
		if err != nil || out.UserID != id.String() || q.id != id || out.Permissions.CanCreateTalkSession != allowed {
			t.Fatalf("out=%+v err=%v", out, err)
		}
	}
	q.out = nil
	if _, err := h.GetCurrentUser(context.Background()); !errors.Is(err, messages.UserNotFoundError) {
		t.Fatal("missing profile fabricated", err)
	}
	q.err = errors.New("database failed")
	if _, err := h.GetCurrentUser(context.Background()); !errors.Is(err, q.err) {
		t.Fatal("db failure hidden", err)
	}
	before := q.calls
	h.authorizationService = currentAuthStub{err: messages.ForbiddenError}
	if _, err := h.GetCurrentUser(context.Background()); !errors.Is(err, messages.ForbiddenError) || q.calls != before {
		t.Fatal("unauthenticated read", err)
	}
}
