package admin_usecase

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

type permissionGrantStub struct {
	granted bool
	calls   int
	err     error
}

func TestPermissionChangeAuditIncludesActorForSuccessAndFailure(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&output, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })
	actor := shared.NewUUID[admin.AdminUser]().String()
	target := shared.NewUUID[user.User]()
	ctx := admin.SetAdminSession(context.Background(), &admin.AdminClaim{Sub: actor})
	permission := &permissionGrantStub{}
	s := NewSetSessionCreationPermission(permission)
	if err := s.Execute(ctx, target, true); err != nil {
		t.Fatal(err)
	}
	permission.err = user.ErrSessionCreationPermissionIneligible
	if err := s.Execute(ctx, target, true); !errors.Is(err, messages.BadRequestError) {
		t.Fatal(err)
	}
	for _, expected := range []string{"session creation permission changed", "session creation permission change failed", actor, target.String(), `"allowed":true`} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("missing audit field %q", expected)
		}
	}
}

func (s *permissionGrantStub) Get(context.Context, shared.UUID[user.User]) (bool, error) {
	panic("permission mutation must not depend on a profile or preflight read")
}
func (s *permissionGrantStub) Set(_ context.Context, _ shared.UUID[user.User], value bool) error {
	s.calls++
	if s.err != nil {
		return s.err
	}
	s.granted = value
	return nil
}

func TestSetSessionCreationPermission(t *testing.T) {
	id := shared.NewUUID[user.User]()
	ctx := admin.SetAdminSession(context.Background(), &admin.AdminClaim{Sub: shared.NewUUID[admin.AdminUser]().String()})
	p := &permissionGrantStub{}
	s := NewSetSessionCreationPermission(p)
	if err := s.Execute(context.Background(), id, true); !errors.Is(err, messages.ForbiddenError) || p.calls != 0 {
		t.Fatal("mutation without admin context was allowed", err)
	}
	for _, allowed := range []bool{true, true, false, false} {
		if err := s.Execute(ctx, id, allowed); err != nil || p.granted != allowed {
			t.Fatalf("allowed=%v err=%v", allowed, err)
		}
	}
	failure := errors.New("permission database unavailable")
	for _, tc := range []struct{ err, want error }{
		{sql.ErrNoRows, messages.UserNotFoundError},
		{user.ErrSessionCreationPermissionIneligible, messages.BadRequestError},
		{failure, failure},
	} {
		p.err = tc.err
		if err := s.Execute(ctx, id, true); !errors.Is(err, tc.want) {
			t.Fatalf("err=%v want=%v", err, tc.want)
		}
	}
}
