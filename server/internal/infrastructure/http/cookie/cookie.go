package cookie

import (
	"net/http"

	"github.com/neko-dream/api/internal/infrastructure/config"
)

const (
	SessionCookieName  = "SessionId"
	SessionMaxAge      = 86400 * 5  // 24 hours
	AuthCookieMaxAge   = 900        // 15 minutes
	GuestSessionMaxAge = 86400 * 30 // 30 days
)

type CookieManager struct {
	config   *config.Config
	domain   string
	secure   bool
	sameSite http.SameSite
}

func NewCookieManager(
	cfg *config.Config,
) CookieManager {
	// dev環境ではフロントエンド(kotohiro.workers.dev)がAPI(*.kotohiro.com)に対して
	// クロスサイトになるため、SameSite=NoneにしないとCookieが送信されない。
	sameSite := http.SameSiteLaxMode
	if cfg.Env == config.DEV {
		sameSite = http.SameSiteNoneMode
	}
	// Local requests are proxied through the Web origin. Host-only cookies
	// work on localhost without depending on a shared domain or hosts entries.
	domain := "." + cfg.DOMAIN
	if cfg.Env == config.LOCAL {
		domain = ""
	}

	return CookieManager{
		config:   cfg,
		domain:   domain,
		secure:   true,
		sameSite: sameSite,
	}
}

func (cm *CookieManager) CreateSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   cm.secure,
		Path:     "/",
		SameSite: cm.sameSite,
		Domain:   cm.domain,
		MaxAge:   SessionMaxAge,
	}
}

func (cm *CookieManager) CreateRevokeCookie() *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   cm.secure,
		Path:     "/",
		SameSite: cm.sameSite,
		Domain:   cm.domain,
		MaxAge:   -1,
	}
}

func (cm *CookieManager) CreateGuestSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   cm.secure,
		Path:     "/",
		SameSite: cm.sameSite,
		Domain:   cm.domain,
		MaxAge:   GuestSessionMaxAge,
	}
}

const AdminSessionCookieName = "AdminSessionId"
const AdminRefreshCookieName = "AdminRefreshToken"
const AdminRefreshMaxAge = 86400 * 30

// AdminSessionMaxAge はアクセストークン（JWT、15分）の Cookie の寿命。
// JWT が切れても Cookie は残し、期限切れの JWT を送らせて 401 → リフレッシュの流れにする
const AdminSessionMaxAge = AdminRefreshMaxAge

func (cm *CookieManager) adminCookie(name, path, value string, age int) *http.Cookie {
	return &http.Cookie{Name: name, Path: path, Value: value, MaxAge: age, HttpOnly: true, Secure: cm.config.Env != config.LOCAL, SameSite: http.SameSiteLaxMode}
}
func (cm *CookieManager) CreateAdminSessionCookie(token string) *http.Cookie {
	return cm.adminCookie(AdminSessionCookieName, "/v1", token, AdminSessionMaxAge)
}
func (cm *CookieManager) CreateAdminRevokeCookie() *http.Cookie {
	return cm.adminCookie(AdminSessionCookieName, "/v1", "", -1)
}
func (cm *CookieManager) CreateAdminRefreshCookie(token string) *http.Cookie {
	return cm.adminCookie(AdminRefreshCookieName, "/v1/admin", token, AdminRefreshMaxAge)
}
func (cm *CookieManager) CreateAdminRefreshRevokeCookie() *http.Cookie {
	return cm.adminCookie(AdminRefreshCookieName, "/v1/admin", "", -1)
}
