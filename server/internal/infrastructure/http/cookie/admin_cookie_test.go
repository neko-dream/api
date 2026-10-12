package cookie

import (
	"github.com/neko-dream/api/internal/infrastructure/config"
	"net/http"
	"testing"
)

func TestAdminCookieScope(t *testing.T) {
	for _, env := range []config.ENV{config.LOCAL, config.DEV, config.PROD} {
		t.Run(string(env), func(t *testing.T) {
			cm := NewCookieManager(&config.Config{Env: env, DOMAIN: "kotohiro.com"})
			for _, c := range []*http.Cookie{cm.CreateAdminSessionCookie("access"), cm.CreateAdminRefreshCookie("refresh"), cm.CreateAdminRevokeCookie(), cm.CreateAdminRefreshRevokeCookie()} {
				if c.Domain != "" || !c.HttpOnly || c.SameSite != http.SameSiteLaxMode || c.Secure != (env != config.LOCAL) {
					t.Fatalf("unsafe admin cookie: %+v", c)
				}
				wantPath := "/v1"
				if c.Name == AdminRefreshCookieName {
					wantPath = "/v1/admin"
				}
				if c.Path != wantPath {
					t.Fatalf("path=%q want=%q", c.Path, wantPath)
				}
			}
			if cm.CreateAdminSessionCookie("access").MaxAge != 900 {
				t.Fatal("access cookie exceeds JWT lifetime")
			}
		})
	}
}
