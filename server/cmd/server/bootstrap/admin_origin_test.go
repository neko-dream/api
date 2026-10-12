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
		t.Run(origin, func(t *testing.T) {
			called := false
			h := b.adminOriginGuard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true; w.WriteHeader(204) }))
			r := httptest.NewRequest("POST", "/v1/admin/refresh", nil)
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
