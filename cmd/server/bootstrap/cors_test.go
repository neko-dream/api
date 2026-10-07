package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neko-dream/api/internal/infrastructure/config"
)

func TestCredentialedCORSUsesExplicitOrigins(t *testing.T) {
	for _, tc := range []struct {
		env          config.ENV
		website, api string
		allowed      []string
	}{
		{config.DEV, "https://develop.kotohiro.com", "https://api.dev.kotohiro.com", []string{"https://develop.kotohiro.com", "https://api.dev.kotohiro.com"}},
		{config.PROD, "https://www.kotohiro.com", "https://api.kotohiro.com", []string{"https://www.kotohiro.com", "https://kotohiro.com", "https://api.kotohiro.com"}},
		{config.LOCAL, "https://localhost:3000", "https://localhost:3000/api", []string{"https://localhost:3000"}},
	} {
		t.Run(string(tc.env), func(t *testing.T) {
			b := &Bootstrap{config: &config.Config{Env: tc.env, WEBSITE_URL: tc.website, EMAIL_VERIFICATION_BASE_URL: tc.api}}
			h := b.apiCORS().Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
			for _, origin := range tc.allowed {
				req := httptest.NewRequest("GET", "/auth/token/info", nil)
				req.Header.Set("Origin", origin)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, req)
				if w.Header().Get("Access-Control-Allow-Origin") != origin || w.Header().Get("Access-Control-Allow-Credentials") != "true" {
					t.Fatal("trusted origin rejected", origin)
				}
			}
			for _, origin := range []string{"https://evil.example", "https://untrusted.kotohiro.com", "https://preview.workers.dev", "null", "http://localhost:3000"} {
				for _, method := range []string{"GET", "OPTIONS"} {
					req := httptest.NewRequest(method, "/auth/token/info", nil)
					req.Header.Set("Origin", origin)
					if method == "OPTIONS" {
						req.Header.Set("Access-Control-Request-Method", "POST")
					}
					w := httptest.NewRecorder()
					h.ServeHTTP(w, req)
					if w.Header().Get("Access-Control-Allow-Origin") != "" {
						t.Fatal("untrusted credentialed CORS allowed", origin)
					}
				}
			}
		})
	}
	b := &Bootstrap{config: &config.Config{Env: "unknown", WEBSITE_URL: "https://evil.example"}}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	b.apiCORS().Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, req)
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("unknown environment CORS failed open")
	}
}

func TestMutationOriginGuardPreventsSimpleFormCSRF(t *testing.T) {
	b := &Bootstrap{config: &config.Config{Env: config.DEV, WEBSITE_URL: "https://develop.kotohiro.com", EMAIL_VERIFICATION_BASE_URL: "https://api.dev.kotohiro.com"}}
	for _, method := range []string{"POST", "PUT", "PATCH", "DELETE"} {
		for _, origin := range []string{"https://evil.example", "https://preview.workers.dev", "null"} {
			called := false
			h := b.browserOriginGuard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
			req := httptest.NewRequest(method, "/users/me", nil)
			req.Header.Set("Origin", origin)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			if w.Code != http.StatusForbidden || called {
				t.Fatal("cross-site mutation accepted", method, origin)
			}
		}
	}
	for _, origin := range []string{"", "https://develop.kotohiro.com", "https://api.dev.kotohiro.com"} {
		called := false
		h := b.browserOriginGuard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set("Origin", origin)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if !called {
			t.Fatal("trusted browser/non-browser request rejected", origin)
		}
	}
}
