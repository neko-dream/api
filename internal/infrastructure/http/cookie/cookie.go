package cookie

import (
	"net/http"

	"github.com/neko-dream/server/internal/infrastructure/config"
)

const (
	SessionCookieName = "SessionId"
	SessionMaxAge     = 86400 * 5 // 24 hours
	AuthCookieMaxAge  = 900       // 15 minutes
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
