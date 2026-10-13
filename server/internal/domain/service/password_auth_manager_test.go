package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/neko-dream/api/internal/domain/messages"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/pkg/hash"
)

type policyTestPasswords struct {
	password_auth.PasswordAuthRepository
	auth             password_auth.PasswordAuth
	creates, updates int
}

func (r *policyTestPasswords) GetPasswordAuthByUserID(context.Context, shared.UUID[user.User]) (password_auth.PasswordAuth, error) {
	return r.auth, nil
}
func (r *policyTestPasswords) CreatePasswordAuth(_ context.Context, _ shared.UUID[user.User], hashed, _ string, _ bool) error {
	r.creates++
	r.auth.PasswordHash = hashed
	return nil
}
func (r *policyTestPasswords) UpdatePasswordAuth(_ context.Context, _ shared.UUID[user.User], hashed, _ string, _ bool) error {
	r.updates++
	r.auth.PasswordHash = hashed
	return nil
}

func TestPasswordAuthManagerSharesRegistrationAndChangePolicy(t *testing.T) {
	cfg := &config.Config{HASH_PEPPER: "local-test-pepper", HASH_ITERATIONS: 4}
	policy, _ := password_auth.NewPasswordPolicy(len(cfg.HASH_PEPPER))
	id := shared.NewUUID[user.User]()
	ctx := context.Background()
	for _, change := range []bool{false, true} {
		for _, tc := range []struct {
			password string
			want     error
		}{
			{"short", messages.PasswordTooShort},
			{strings.Repeat("😀", 11), messages.PasswordTooShort},
			{strings.Repeat("a", policy.MaxBytes+1), messages.PasswordTooLong},
			{strings.Repeat("😀", 14), messages.PasswordTooLong},
			{strings.Repeat("a", policy.MaxBytes), nil},
			{strings.Repeat("あ", 12), nil},
		} {
			repo := &policyTestPasswords{}
			manager := NewPasswordAuthManager(repo, cfg)
			var err error
			if change {
				err = manager.UpdatePassword(ctx, id, tc.password)
			} else {
				err = manager.RegisterPassword(ctx, id, tc.password, false)
			}
			if !errors.Is(err, tc.want) {
				t.Fatalf("change %v: %v want %v", change, err, tc.want)
			}
			writes := repo.creates + repo.updates
			if tc.want != nil && writes != 0 {
				t.Fatal("invalid password written")
			}
			if tc.want == nil {
				valid, err := manager.VerifyPassword(ctx, id, tc.password)
				if writes != 1 || err != nil || !valid {
					t.Fatal("accepted password was not usable", err)
				}
			}
		}
	}
}

func TestPasswordAuthManagerLoginKeepsLegacyPasswordCompatibility(t *testing.T) {
	cfg := &config.Config{HASH_PEPPER: "local-test-pepper", HASH_ITERATIONS: 4}
	legacy := "old-pass"
	hashed, err := hash.BcryptHash(legacy, cfg.HASH_PEPPER, cfg.HASH_ITERATIONS)
	if err != nil {
		t.Fatal(err)
	}
	repo := &policyTestPasswords{auth: password_auth.PasswordAuth{PasswordHash: hashed}}
	manager := NewPasswordAuthManager(repo, cfg)
	id := shared.NewUUID[user.User]()
	valid, err := manager.VerifyPassword(context.Background(), id, legacy)
	if err != nil || !valid || repo.creates+repo.updates != 0 {
		t.Fatal("new policy blocked an existing password", err)
	}
	valid, err = manager.VerifyPassword(context.Background(), id, "wrong")
	if err != nil || valid {
		t.Fatal("wrong legacy password accepted", err)
	}
}
