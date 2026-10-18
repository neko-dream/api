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
			// アクセストークンの Cookie は、JWT（15分）が切れても送られる必要がある。
			// 期限切れの JWT は 401 になり、管理画面がリフレッシュする。Cookie ごと消えると 403 になり、
			// リフレッシュされずにログイン画面に戻される
			if got := cm.CreateAdminSessionCookie("access").MaxAge; got != cm.CreateAdminRefreshCookie("refresh").MaxAge {
				t.Fatalf("access cookie MaxAge=%d must match refresh cookie lifetime", got)
			}
		})
	}
}
