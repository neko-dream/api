package admin_usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	adminauth "github.com/neko-dream/api/internal/infrastructure/auth/admin"
	"github.com/neko-dream/api/internal/infrastructure/config"
)

type authTestUsers struct {
	admin.AdminUserRepository
	user *admin.AdminUser
}

func (r *authTestUsers) FindByEmail(context.Context, string) (*admin.AdminUser, error) {
	return r.user, nil
}
func (r *authTestUsers) FindByID(context.Context, shared.UUID[admin.AdminUser]) (*admin.AdminUser, error) {
	return r.user, nil
}

type authTestVerifier struct{ err error }

func (v authTestVerifier) Verify(context.Context, string) (*admin.GoogleClaims, error) {
	return &admin.GoogleClaims{Email: "operator@example.com"}, v.err
}

type authTestSessions struct {
	admin.AdminSessionRepository
	row *admin.AdminSession
	// beforeRotate は Rotate の直前に呼ばれる。別タブが先にローテーションした状況を再現する
	beforeRotate func()
}

func (r *authTestSessions) Create(_ context.Context, row admin.AdminSession) error {
	r.row = &row
	return nil
}
func (r *authTestSessions) matches(hash string) bool {
	return r.row != nil && (r.row.TokenHash == hash || (r.row.PreviousTokenHash != nil && *r.row.PreviousTokenHash == hash))
}
func (r *authTestSessions) FindByTokenHash(_ context.Context, hash string) (*admin.AdminSession, error) {
	if r.matches(hash) {
		copied := *r.row
		return &copied, nil
	}
	return nil, nil
}
func (r *authTestSessions) Rotate(_ context.Context, id shared.UUID[admin.AdminSession], oldHash, hash string, expires, rotatedAt time.Time) (bool, error) {
	if hook := r.beforeRotate; hook != nil {
		r.beforeRotate = nil
		hook()
	}
	if r.row == nil || r.row.AdminSessionID != id {
		return false, errors.New("unknown session")
	}
	if r.row.TokenHash != oldHash {
		return false, nil
	}
	previous := r.row.TokenHash
	r.row.PreviousTokenHash, r.row.RotatedAt = &previous, &rotatedAt
	r.row.TokenHash, r.row.ExpiresAt = hash, expires
	return true, nil
}
func (r *authTestSessions) DeleteByTokenHash(_ context.Context, hash string) error {
	if r.matches(hash) {
		r.row = nil
	}
	return nil
}

func TestAdminLoginRefreshAndLogout(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{AdminJWTSecret: "local-test-only-secret-at-least-32-bytes", HASH_PEPPER: "local-test-only-pepper"}
	users := &authTestUsers{user: &admin.AdminUser{AdminUserID: shared.NewUUID[admin.AdminUser](), Email: "operator@example.com", DisplayName: "Operator"}}
	sessions := &authTestSessions{}
	tokens := adminauth.NewAdminTokenManager(cfg)
	login := NewAdminLogin(users, sessions, tokens, authTestVerifier{}, cfg)
	refresh := NewAdminRefresh(users, sessions, tokens, cfg)
	logout := NewAdminLogout(sessions, cfg)
	out, err := login.Execute(ctx, AdminLoginInput{IDToken: "verified-by-test-double"})
	if err != nil {
		t.Fatal(err)
	}
	claim, err := tokens.Parse(ctx, out.AccessToken)
	if err != nil || claim.Sub != users.user.AdminUserID.String() {
		t.Fatalf("access token: %v %v", claim, err)
	}
	if sessions.row == nil || sessions.row.TokenHash == out.RefreshToken {
		t.Fatal("refresh token missing or stored in plaintext")
	}
	updated, err := refresh.Execute(ctx, out.RefreshToken)
	if err != nil || updated.RefreshToken == "" || updated.RefreshToken == out.RefreshToken {
		t.Fatalf("refresh rotation: %v", err)
	}
	// 猶予期間内の直前のトークンは、ローテーションせずアクセストークンだけを発行する
	concurrent, err := refresh.Execute(ctx, out.RefreshToken)
	if err != nil || concurrent.AccessToken == "" || concurrent.RefreshToken != "" {
		t.Fatalf("previous refresh within grace period: %+v %v", concurrent, err)
	}
	afterGrace := clock.SetNow(ctx, time.Now().Add(admin.RotationGracePeriod+time.Second))
	if _, err := refresh.Execute(afterGrace, out.RefreshToken); !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("old refresh accepted after grace period: %v", err)
	}
	updated, err = refresh.Execute(ctx, updated.RefreshToken)
	if err != nil || updated.RefreshToken == "" {
		t.Fatalf("current refresh was invalidated by the previous one: %v", err)
	}
	if err := logout.Execute(ctx, updated.RefreshToken); err != nil {
		t.Fatal(err)
	}
	if _, err := refresh.Execute(ctx, updated.RefreshToken); !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("logout did not revoke refresh: %v", err)
	}
	if _, err := refresh.Execute(ctx, ""); !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("empty refresh accepted: %v", err)
	}
	out, err = login.Execute(ctx, AdminLoginInput{IDToken: "verified-by-test-double"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := refresh.Execute(clock.SetNow(ctx, sessions.row.ExpiresAt.Add(time.Second)), out.RefreshToken); !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("expired refresh accepted: %v", err)
	}
}

func TestAdminLoginRejectsUnverifiedAndUnlistedAccounts(t *testing.T) {
	cfg := &config.Config{AdminJWTSecret: "local-test-only-secret-at-least-32-bytes", HASH_PEPPER: "local-test-only-pepper"}
	for _, verificationError := range []error{nil, errors.New("invalid Google token")} {
		sessions := &authTestSessions{}
		login := NewAdminLogin(&authTestUsers{}, sessions, adminauth.NewAdminTokenManager(cfg), authTestVerifier{verificationError}, cfg)
		if _, err := login.Execute(context.Background(), AdminLoginInput{}); !errors.Is(err, ErrInvalidAdminCredentials) {
			t.Fatalf("invalid/unlisted account: %v", err)
		}
		if sessions.row != nil {
			t.Fatal("rejected login created session")
		}
	}
}

func TestAdminRefreshConcurrentRotationSucceedsWithinGracePeriod(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{AdminJWTSecret: "local-test-only-secret-at-least-32-bytes", HASH_PEPPER: "local-test-only-pepper"}
	users := &authTestUsers{user: &admin.AdminUser{AdminUserID: shared.NewUUID[admin.AdminUser](), Email: "operator@example.com", DisplayName: "Operator"}}
	sessions := &authTestSessions{}
	tokens := adminauth.NewAdminTokenManager(cfg)
	login := NewAdminLogin(users, sessions, tokens, authTestVerifier{}, cfg)
	refresh := NewAdminRefresh(users, sessions, tokens, cfg)
	out, err := login.Execute(ctx, AdminLoginInput{IDToken: "verified-by-test-double"})
	if err != nil {
		t.Fatal(err)
	}

	// 同じリフレッシュトークンで、別タブが先にローテーションを終えた
	sessions.beforeRotate = func() {
		if _, err := refresh.Execute(ctx, out.RefreshToken); err != nil {
			t.Fatal(err)
		}
	}
	got, err := refresh.Execute(ctx, out.RefreshToken)
	if err != nil || got.AccessToken == "" || got.RefreshToken != "" {
		t.Fatalf("concurrent refresh: %+v %v", got, err)
	}
}

// リフレッシュトークンは Cookie に入れる。ogen は Set-Cookie の値を URL エンコードするので、
// URL で使える文字だけにして、ブラウザに保存される値が変わらないようにする。
func TestAdminRefreshTokensAreURLSafe(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{AdminJWTSecret: "local-test-only-secret-at-least-32-bytes", HASH_PEPPER: "local-test-only-pepper"}
	users := &authTestUsers{user: &admin.AdminUser{AdminUserID: shared.NewUUID[admin.AdminUser](), Email: "operator@example.com", DisplayName: "Operator"}}
	sessions := &authTestSessions{}
	tokens := adminauth.NewAdminTokenManager(cfg)
	login := NewAdminLogin(users, sessions, tokens, authTestVerifier{}, cfg)
	refresh := NewAdminRefresh(users, sessions, tokens, cfg)

	for i := 0; i < 20; i++ {
		out, err := login.Execute(ctx, AdminLoginInput{IDToken: "verified-by-test-double"})
		if err != nil {
			t.Fatal(err)
		}
		rotated, err := refresh.Execute(ctx, out.RefreshToken)
		if err != nil {
			t.Fatal(err)
		}
		for _, token := range []string{out.RefreshToken, rotated.RefreshToken} {
			if len(token) < 64 || strings.Trim(token, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_") != "" {
				t.Fatalf("refresh token is not URL-safe: %q", token)
			}
		}
	}
}
