package cookie

import (
	"net/http"
	"testing"

	"github.com/neko-dream/api/internal/infrastructure/config"
)

func TestCookieEnvironmentPolicy(t *testing.T) {
	for _, tt := range []struct {
		name     string
		env      config.ENV
		domain   string
		sameSite http.SameSite
	}{
		{"local", config.LOCAL, "", http.SameSiteLaxMode},
		{"development", config.DEV, ".kotohiro.com", http.SameSiteNoneMode},
		{"production", config.PROD, ".kotohiro.com", http.SameSiteLaxMode},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewCookieManager(&config.Config{Env: tt.env, DOMAIN: "kotohiro.com"})
			for _, c := range []*http.Cookie{cm.CreateSessionCookie("token"), cm.CreateGuestSessionCookie("token"), cm.CreateRevokeCookie()} {
				if c.Domain != tt.domain || c.SameSite != tt.sameSite || !c.Secure || !c.HttpOnly || c.Path != "/" {
					t.Fatalf("unexpected cookie policy: %s", c.String())
				}
			}
		})
	}
}
