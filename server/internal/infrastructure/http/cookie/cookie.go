package cookie

import (
	"net/http"

	"github.com/neko-dream/api/internal/infrastructure/config"
)

const (
	SessionCookieName      = "SessionId"
	AdminSessionCookieName = "AdminSessionId"
	AdminRefreshCookieName = "AdminRefreshToken"
	SessionMaxAge          = 86400 * 5  // 24 hours
	AuthCookieMaxAge       = 900        // 15 minutes
	GuestSessionMaxAge     = 86400 * 30 // 30 days
	AdminSessionMaxAge     = 86400 * 7  // 7 days
	AdminRefreshMaxAge     = 86400 * 30 // 30 days
)

type CookieManager struct {
	config   *config.Config
	secure   bool
	sameSite http.SameSite
}

func NewCookieManager(
	config *config.Config,
) CookieManager {
	return CookieManager{
		config:   config,
		secure:   true,
		sameSite: http.SameSiteLaxMode,
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
		Domain:   "." + cm.config.DOMAIN,
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
		SameSite: http.SameSiteLaxMode,
		Domain:   "." + cm.config.DOMAIN,
		MaxAge:   -1,
	}
}

func (cm *CookieManager) CreateAdminSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     AdminSessionCookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   cm.secure,
		Path:     "/",
		SameSite: cm.sameSite,
		Domain:   "." + cm.config.DOMAIN,
		MaxAge:   AdminSessionMaxAge,
	}
}

func (cm *CookieManager) CreateAdminRevokeCookie() *http.Cookie {
	return &http.Cookie{
		Name:     AdminSessionCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   cm.secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Domain:   "." + cm.config.DOMAIN,
		MaxAge:   -1,
	}
}

func (cm *CookieManager) CreateAdminRefreshCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     AdminRefreshCookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   cm.secure,
		Path:     "/",
		SameSite: cm.sameSite,
		Domain:   "." + cm.config.DOMAIN,
		MaxAge:   AdminRefreshMaxAge,
	}
}

func (cm *CookieManager) CreateAdminRefreshRevokeCookie() *http.Cookie {
	return &http.Cookie{
		Name:     AdminRefreshCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   cm.secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Domain:   "." + cm.config.DOMAIN,
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
		Domain:   "." + cm.config.DOMAIN,
		MaxAge:   GuestSessionMaxAge,
	}
}
