package bootstrap

import (
	"github.com/neko-dream/api/internal/infrastructure/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminOriginGuard(t *testing.T) {
	b := &Bootstrap{config: &config.Config{AdminAllowedOrigins: []string{"https://admin.example.com"}}}
	for _, origin := range []string{"", "null", "https://evil.example.com", "https://admin.example.com"} {
		for _, endpoint := range []struct{ method, path string }{{"POST", "/v1/admin/refresh"}, {"PUT", "/v1/manage/users/user-id/permissions"}} {
			t.Run(origin+endpoint.path, func(t *testing.T) {
				called := false
				h := b.adminOriginGuard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true; w.WriteHeader(204) }))
				r := httptest.NewRequest(endpoint.method, endpoint.path, nil)
				if origin != "" {
					r.Header.Set("Origin", origin)
				}
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				if called != (origin == "https://admin.example.com") {
					t.Fatalf("origin %q reached handler=%v status=%d", origin, called, w.Code)
				}
			})
		}
	}
}

func TestAdminPermissionRouteUsesDedicatedOriginGuard(t *testing.T) {
	b := &Bootstrap{config: &config.Config{Env: config.DEV, WEBSITE_URL: "https://web.example.com", AdminAllowedOrigins: []string{"https://admin.example.com"}}}
	for _, path := range []string{"/v1/manage/users/user-id/permissions", "/v1/admin/refresh"} {
		for _, origin := range []string{"", "https://web.example.com", "https://evil.example.com", "https://admin.example.com"} {
			called := false
			h := b.wrapAPIOrigins(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { called = true; w.WriteHeader(http.StatusNoContent) }))
			r := httptest.NewRequest(http.MethodPut, path, nil)
			if origin != "" {
				r.Header.Set("Origin", origin)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if called != (origin == "https://admin.example.com") {
				t.Fatalf("path=%s origin=%s status=%d", path, origin, w.Code)
			}
		}
	}
	called := false
	h := b.wrapAPIOrigins(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { called = true; w.WriteHeader(http.StatusNoContent) }))
	r := httptest.NewRequest(http.MethodPost, "/talksessions", nil)
	r.Header.Set("Origin", "https://web.example.com")
	h.ServeHTTP(httptest.NewRecorder(), r)
	if !called {
		t.Fatal("ordinary API must keep the Web origin guard")
	}
}
