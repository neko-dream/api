package admin_usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

type permissionCurrentStub struct {
	output *userquery.CurrentOutput
	err    error
}

func (s *permissionCurrentStub) FindByID(context.Context, shared.UUID[user.User]) (*userquery.CurrentOutput, error) {
	return s.output, s.err
}

type permissionGrantStub struct {
	granted bool
	calls   int
	err     error
}

func (s *permissionGrantStub) Get(context.Context, shared.UUID[user.User]) (bool, error) {
	return s.granted, s.err
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
	q := &permissionCurrentStub{output: &userquery.CurrentOutput{}}
	p := &permissionGrantStub{}
	s := NewSetSessionCreationPermission(q, p)
	for _, allowed := range []bool{true, true, false, false} {
		if err := s.Execute(context.Background(), id, allowed); err != nil || p.granted != allowed {
			t.Fatalf("allowed=%v err=%v", allowed, err)
		}
	}
	q.output.IsGuest = true
	before := p.calls
	if err := s.Execute(context.Background(), id, true); !errors.Is(err, messages.BadRequestError) || p.calls != before {
		t.Fatal("guest grant accepted", err)
	}
	if err := s.Execute(context.Background(), id, false); err != nil {
		t.Fatal("guest revoke rejected", err)
	}
	q.output.IsGuest = false
	now := time.Now()
	q.output.User.WithdrawalDate = &now
	if err := s.Execute(context.Background(), id, true); !errors.Is(err, messages.BadRequestError) {
		t.Fatal("withdrawn grant accepted", err)
	}
	q.output = nil
	if err := s.Execute(context.Background(), id, true); !errors.Is(err, messages.UserNotFoundError) {
		t.Fatal("missing user accepted", err)
	}
	q.err = errors.New("db failed")
	if err := s.Execute(context.Background(), id, true); !errors.Is(err, q.err) {
		t.Fatal("db error hidden", err)
	}
	q.err = nil
	q.output = &userquery.CurrentOutput{}
	p.err = errors.New("write failed")
	if err := s.Execute(context.Background(), id, true); !errors.Is(err, p.err) {
		t.Fatal("write error hidden", err)
	}
}
