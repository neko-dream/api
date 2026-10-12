package bootstrap

import (
	"context"
	"errors"
	"testing"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"go.uber.org/dig"
)

type initialAdminRepo struct {
	admin.AdminUserRepository
	user       *admin.AdminUser
	lookupErr  error
	createErr  error
	concurrent bool
	creates    int
}

func (r *initialAdminRepo) FindByEmail(_ context.Context, email string) (*admin.AdminUser, error) {
	if r.user != nil && r.user.Email != email {
		return nil, nil
	}
	return r.user, r.lookupErr
}

func (r *initialAdminRepo) Create(_ context.Context, user admin.AdminUser) error {
	r.creates++
	if r.createErr == nil || r.concurrent {
		r.user = &user
	}
	return r.createErr
}

func TestInitialAdmin(t *testing.T) {
	failure := errors.New("database unavailable")
	for _, tt := range []struct {
		name        string
		email       string
		display     string
		repo        initialAdminRepo
		wantErr     bool
		wantCreates int
	}{
		{name: "unset does not register anyone"},
		{name: "register configured email", email: "operator@example.com", display: "Operator", wantCreates: 1},
		{name: "default name", email: "operator@example.com", wantCreates: 1},
		{name: "existing user unchanged", email: "operator@example.com", repo: initialAdminRepo{user: &admin.AdminUser{Email: "operator@example.com", DisplayName: "Existing"}}},
		{name: "lookup failure fails startup", email: "operator@example.com", repo: initialAdminRepo{lookupErr: failure}, wantErr: true},
		{name: "create failure fails startup", email: "operator@example.com", repo: initialAdminRepo{createErr: failure}, wantErr: true, wantCreates: 1},
		{name: "concurrent registration", email: "operator@example.com", repo: initialAdminRepo{createErr: failure, concurrent: true}, wantCreates: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			container := dig.New()
			if err := container.Provide(func() admin.AdminUserRepository { return &tt.repo }); err != nil {
				t.Fatal(err)
			}
			b := &Bootstrap{container: container, config: &config.Config{InitialAdminEmail: tt.email, InitialAdminName: tt.display}}
			err := b.bootstrapInitialAdmin(context.Background())
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v", err)
			}
			if tt.repo.creates != tt.wantCreates {
				t.Fatalf("creates = %d", tt.repo.creates)
			}
			if tt.wantCreates == 1 && !tt.wantErr {
				wantName := tt.display
				if wantName == "" {
					wantName = tt.email
				}
				if tt.repo.user.Email != tt.email || tt.repo.user.DisplayName != wantName {
					t.Fatalf("unexpected user: %+v", tt.repo.user)
				}
				if err := b.bootstrapInitialAdmin(context.Background()); err != nil {
					t.Fatal(err)
				}
				if tt.repo.creates != 1 {
					t.Fatal("restart created duplicate admin")
				}
			}
		})
	}
}
