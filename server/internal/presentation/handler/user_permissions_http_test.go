package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	adminauth "github.com/neko-dream/api/internal/infrastructure/auth/admin"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/presentation/oas"
	httputils "github.com/neko-dream/api/pkg/http"
)

type httpPermissionChange struct {
	calls   int
	allowed bool
}

func (s *httpPermissionChange) Execute(_ context.Context, _ shared.UUID[user.User], value bool) error {
	s.calls++
	s.allowed = value
	return nil
}

func permissionRequest(srv http.Handler, method, path, cookie, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	if cookie != "" {
		r.Header.Set("Cookie", cookie)
	}
	w := httptest.NewRecorder()
	r = r.WithContext(httputils.WithHTTPResReqContext(r.Context(), r, w))
	srv.ServeHTTP(w, r)
	return w
}

func TestPermissionMutationRequiresDedicatedAdminCookie(t *testing.T) {
	mgr := adminauth.NewAdminTokenManager(&config.Config{AdminJWTSecret: "operator-permissions-http-test-secret-32bytes"})
	token, err := mgr.Generate(context.Background(), admin.AdminUser{AdminUserID: shared.NewUUID[admin.AdminUser]()})
	if err != nil {
		t.Fatal(err)
	}
	change := &httpPermissionChange{}
	srv, err := oas.NewServer(&handlers{ManageHandler: &manageHandler{setCreationPermission: change}}, NewSecurityHandler(nil, nil, mgr, adminSecurityUsers{}), oas.WithErrorHandler(CustomErrorHandler))
	if err != nil {
		t.Fatal(err)
	}
	path := "/v1/manage/users/" + shared.NewUUID[user.User]().String() + "/permissions"
	for _, tc := range []struct {
		cookie, body string
		status       int
	}{
		{"", `{"canCreateTalkSession":true}`, 403},
		{"SessionId=" + token, `{"canCreateTalkSession":true}`, 403},
		{"SessionId=organization-session", `{"canCreateTalkSession":true}`, 403},
		{"AdminSessionId=invalid", `{"canCreateTalkSession":true}`, 403},
		{"AdminSessionId=" + token, `{}`, 400},
		{"AdminSessionId=" + token, `{"canCreateTalkSession":true}`, 200},
		{"AdminSessionId=" + token, `{"canCreateTalkSession":false}`, 200},
	} {
		before := change.calls
		w := permissionRequest(srv, "PUT", path, tc.cookie, tc.body)
		if w.Code != tc.status {
			t.Fatalf("got %d want %d: %s", w.Code, tc.status, w.Body)
		}
		if tc.status != 200 && change.calls != before {
			t.Fatal("unauthorized or invalid permission mutation reached usecase")
		}
	}
	if change.calls != 2 || change.allowed {
		t.Fatal("grant/revoke route mismatch")
	}
}

type meSessionToken struct {
	session.TokenManager
	claim *session.Claim
}

func (s meSessionToken) Parse(_ context.Context, token string) (*session.Claim, error) {
	if token != "test-user-session" {
		return nil, errors.New("invalid")
	}
	return s.claim, nil
}

type meSessions struct {
	session.SessionRepository
	value *session.Session
}

func (s meSessions) FindBySessionID(context.Context, shared.UUID[session.Session]) (*session.Session, error) {
	return s.value, nil
}

func TestCurrentUserHTTPUsesUserCookieAndNeverReturnsTokenOrOrganization(t *testing.T) {
	id := shared.NewUUID[user.User]()
	sid := shared.NewUUID[session.Session]()
	claim := &session.Claim{Sub: id.String(), Jti: sid.String(), IsRegistered: true, Exp: time.Now().Add(time.Hour).Unix()}
	q := &currentQueryStub{out: &userquery.CurrentOutput{CanCreateTalkSession: true}}
	srv, err := oas.NewServer(&handlers{UserHandler: &userHandler{authorizationService: currentAuthStub{id: id}, currentUser: q}}, NewSecurityHandler(meSessionToken{claim: claim}, meSessions{value: session.NewSession(sid, id, shared.ProviderGoogle, session.SESSION_ACTIVE, time.Now().Add(time.Hour), time.Now())}, nil, nil), oas.WithErrorHandler(CustomErrorHandler))
	if err != nil {
		t.Fatal(err)
	}
	for _, cookie := range []string{"", "AdminSessionId=test-admin", "SessionId=invalid"} {
		w := permissionRequest(srv, "GET", "/users/me", cookie, "")
		if w.Code != 403 {
			t.Fatalf("unauthenticated current-user status=%d", w.Code)
		}
	}
	for _, value := range []bool{true, false} {
		q.out.CanCreateTalkSession = value
		w := permissionRequest(srv, "GET", "/users/me", "SessionId=test-user-session", "")
		if w.Code != 200 || !strings.Contains(w.Body.String(), `"canCreateTalkSession":`+map[bool]string{true: "true", false: "false"}[value]) {
			t.Fatalf("me status=%d body=%s", w.Code, w.Body)
		}
		if w.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("current-user response must not be cached")
		}
		for _, private := range []string{"organization", "orgType", "subject", "email", "demographics", "jti", "exp", "iat"} {
			if strings.Contains(w.Body.String(), `"`+private+`"`) {
				t.Fatal("current-user leaked token/private field", private)
			}
		}
	}
}
