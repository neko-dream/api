package cookie

import (
	"net/http"

	"github.com/neko-dream/server/internal/infrastructure/config"
)

const (
	SessionCookieName  = "SessionId"
	StateCookieName    = "state"
	RedirectCookieName = "redirect_url"
	SessionMaxAge      = 86400 // 24 hours
	AuthCookieMaxAge   = 300   // 5 minutes
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
		Domain:   cm.config.DOMAIN,
		MaxAge:   SessionMaxAge,
	}
}

func (cm *CookieManager) CreateAuthCookies(state, redirectURL string) []*http.Cookie {
	return []*http.Cookie{
		{
			Name:     StateCookieName,
			Value:    state,
			HttpOnly: true,
			Secure:   cm.secure,
			Path:     "/",
			SameSite: cm.sameSite,
			Domain:   cm.config.DOMAIN,
			MaxAge:   AuthCookieMaxAge,
		},
		{
			Name:     RedirectCookieName,
			Value:    redirectURL,
			HttpOnly: true,
			Secure:   cm.secure,
			Path:     "/",
			SameSite: cm.sameSite,
			Domain:   cm.config.DOMAIN,
			MaxAge:   AuthCookieMaxAge,
		},
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
		Domain:   cm.config.DOMAIN,
		MaxAge:   -1,
	}
}
