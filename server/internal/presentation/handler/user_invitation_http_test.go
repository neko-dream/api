package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/neko-dream/api/internal/application/usecase/admin_usecase"
	"github.com/neko-dream/api/internal/application/usecase/auth_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	adminauth "github.com/neko-dream/api/internal/infrastructure/auth/admin"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/presentation/oas"
	httputils "github.com/neko-dream/api/pkg/http"
)

func TestPasswordPolicyHTTPIsPublicAndDoesNotExposeHashingSecrets(t *testing.T) {
	for _, pepper := range []string{"", "local-test-pepper", strings.Repeat("あ", 6)} {
		cfg := &config.Config{Env: config.LOCAL, HASH_PEPPER: pepper}
		srv, err := oas.NewServer(&handlers{AuthHandler: &authHandler{cfg: cfg}}, NewSecurityHandler(nil, nil, nil, nil), oas.WithErrorHandler(CustomErrorHandler))
		if err != nil {
			t.Fatal(err)
		}
		r := httptest.NewRequest("GET", "/auth/password/policy", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(httputils.WithHTTPResReqContext(r.Context(), r, w))
		srv.ServeHTTP(w, r)
		if w.Code != 200 || w.Header().Get("Set-Cookie") != "" || w.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("unsafe policy response", w.Code)
		}
		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		policy, _ := password_auth.NewPasswordPolicy(len(pepper))
		if len(body) != 4 || body["minLength"] != float64(policy.MinLength) || body["maxBytes"] != float64(policy.MaxBytes) ||
			body["forbidControlCharacters"] != policy.ForbidControlCharacters || body["forbidSurroundingWhitespace"] != policy.ForbidSurroundingWhitespace {
			t.Fatal("public policy does not match validation", body)
		}
	}
	bad := &config.Config{HASH_PEPPER: strings.Repeat("x", 61)}
	srv, err := oas.NewServer(&handlers{AuthHandler: &authHandler{cfg: bad}}, NewSecurityHandler(nil, nil, nil, nil), oas.WithErrorHandler(CustomErrorHandler))
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", "/auth/password/policy", nil)
	w := httptest.NewRecorder()
	r = r.WithContext(httputils.WithHTTPResReqContext(r.Context(), r, w))
	srv.ServeHTTP(w, r)
	if w.Code != 500 || strings.Contains(w.Body.String(), bad.HASH_PEPPER) {
		t.Fatal("bad hashing configuration exposed or accepted")
	}
}

type httpUserInvitation struct {
	calls int
	err   error
}

func (s *httpUserInvitation) Create(context.Context, admin_usecase.InviteUserInput) (*admin_usecase.InviteUserOutput, error) {
	s.calls++
	if s.err != nil {
		return nil, s.err
	}
	return &admin_usecase.InviteUserOutput{UserID: shared.NewUUID[user.User](), InvitationSent: false}, nil
}
func (s *httpUserInvitation) Resend(ctx context.Context, id shared.UUID[user.User]) (*admin_usecase.InviteUserOutput, error) {
	return s.Create(ctx, admin_usecase.InviteUserInput{})
}

func TestUserInvitationHTTPAdminBoundary(t *testing.T) {
	mgr := adminauth.NewAdminTokenManager(&config.Config{AdminJWTSecret: "operator-invitation-http-test-secret-32bytes"})
	token, err := mgr.Generate(context.Background(), admin.AdminUser{AdminUserID: shared.NewUUID[admin.AdminUser]()})
	if err != nil {
		t.Fatal(err)
	}
	invitation := &httpUserInvitation{}
	srv, err := oas.NewServer(&handlers{ManageHandler: &manageHandler{userInvitations: invitation}}, NewSecurityHandler(nil, nil, mgr, adminSecurityUsers{}), oas.WithErrorHandler(CustomErrorHandler))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/v1/manage/users", "/v1/manage/users/" + shared.NewUUID[user.User]().String() + "/invitation"} {
		for _, cookie := range []string{"", "SessionId=" + token, "SessionId=organization-session", "AdminSessionId=invalid", "AdminSessionId=" + token} {
			before := invitation.calls
			w := permissionRequest(srv, "POST", path, cookie, `{"displayName":"表示名","email":"user@example.com"}`)
			want := 403
			if cookie == "AdminSessionId="+token {
				want = 200
				if path == "/v1/manage/users" {
					want = 201
				}
			}
			if w.Code != want {
				t.Fatalf("%s: %d want %d: %s", path, w.Code, want, w.Body)
			}
			if want == 403 && invitation.calls != before {
				t.Fatal("non-admin provisioned a user")
			}
			if want != 403 {
				if !strings.Contains(w.Body.String(), `"invitationSent":false`) {
					t.Fatal("delivery failure status missing")
				}
				for _, secret := range []string{"password", "token", "organization", "email"} {
					if strings.Contains(w.Body.String(), `"`+secret+`"`) {
						t.Fatal("creation response leaked credentials")
					}
				}
			}
		}
	}
	for _, body := range []string{`{}`, `{"displayName":"","email":"user@example.com"}`, `{"displayName":"名前","email":""}`} {
		before := invitation.calls
		w := permissionRequest(srv, "POST", "/v1/manage/users", "AdminSessionId="+token, body)
		if w.Code != 400 || invitation.calls != before {
			t.Fatal("invalid create reached usecase", w.Body.String())
		}
	}
	for _, failure := range []error{messages.ManagedUserEmailExists, messages.PasswordSetupRateLimit} {
		invitation.err = failure
		w := permissionRequest(srv, "POST", "/v1/manage/users", "AdminSessionId="+token, `{"displayName":"名前","email":"user@example.com"}`)
		var api *messages.APIError
		errors.As(failure, &api)
		if w.Code != api.StatusCode || !strings.Contains(w.Body.String(), api.Code) {
			t.Fatal("actionable API error lost")
		}
	}
}

type httpPasswordSetup struct{ calls int }

func (s *httpPasswordSetup) Execute(_ context.Context, input auth_usecase.SetupPasswordInput) error {
	policy, _ := password_auth.NewPasswordPolicy(0)
	if err := policy.Validate(input.Password); err != nil {
		return err
	}
	s.calls++
	return nil
}

func TestPasswordSetupHTTPNoLoginAndNoGetConsumption(t *testing.T) {
	setup := &httpPasswordSetup{}
	srv, err := oas.NewServer(&handlers{AuthHandler: &authHandler{setupPassword: setup, cfg: &config.Config{Env: config.LOCAL, WEBSITE_URL: "https://example.com"}}}, NewSecurityHandler(nil, nil, nil, nil), oas.WithErrorHandler(CustomErrorHandler))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		method, origin, body string
		status               int
	}{
		{"POST", "https://example.com", `{"token":"` + strings.Repeat("A", 43) + `=","password":"new-password-for-setup"}`, 204},
		{"POST", "https://evil.example", `{"token":"` + strings.Repeat("A", 43) + `=","password":"new-password-for-setup"}`, 403},
		{"POST", "null", `{"token":"` + strings.Repeat("A", 43) + `=","password":"new-password-for-setup"}`, 403},
		{"POST", "https://example.com", `{}`, 400},
		{"POST", "https://example.com", `{"token":"short","password":"new-password-for-setup"}`, 400},
		{"POST", "https://example.com", `{"token":"` + strings.Repeat("A", 43) + `=","password":"short"}`, 400},
		{"GET", "https://example.com", "", 405},
	} {
		before := setup.calls
		r := httptest.NewRequest(tc.method, "/auth/password/setup", strings.NewReader(tc.body))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Origin", tc.origin)
		w := httptest.NewRecorder()
		r = r.WithContext(httputils.WithHTTPResReqContext(r.Context(), r, w))
		srv.ServeHTTP(w, r)
		if w.Code != tc.status {
			t.Fatalf("%s %s: %d want %d: %s", tc.method, tc.origin, w.Code, tc.status, w.Body)
		}
		if tc.status != 204 && setup.calls != before {
			t.Fatal("untrusted, invalid or GET setup consumed a token")
		}
		if w.Header().Get("Set-Cookie") != "" {
			t.Fatal("password setup must not log in")
		}
		if tc.status == http.StatusNoContent && (w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Referrer-Policy") != "no-referrer") {
			t.Fatal("setup missing security headers")
		}
	}
}
