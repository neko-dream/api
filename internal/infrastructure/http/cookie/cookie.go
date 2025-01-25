package cookie

import (
	"net/http"

	"github.com/neko-dream/server/internal/infrastructure/config"
)

type CookieManager struct {
	config   *config.Config
	secure   bool
	someSite http.SameSite
}

func NewCookieManager(
	config *config.Config,
) *CookieManager {
	return &CookieManager{
		config:   config,
		secure:   true,
		someSite: http.SameSiteLaxMode,
	}
}

func (cm *CookieManager) CreateSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "SessionId",
		Value:    token,
		HttpOnly: true,
		Secure:   cm.secure,
		Path:     "/",
		SameSite: cm.someSite,
		Domain:   cm.config.DOMAIN,
		MaxAge:   60 * 60 * 24,
	}
}

func (cm *CookieManager) CreateAuthCookies(state, redirectURL string) []*http.Cookie {
	return []*http.Cookie{
		{
			Name:     "state",
			Value:    state,
			HttpOnly: true,
			Secure:   cm.secure,
			Path:     "/",
			SameSite: cm.someSite,
			Domain:   cm.config.DOMAIN,
		},
		{
			Name:     "redirect_url",
			Value:    redirectURL,
			HttpOnly: true,
			Secure:   cm.secure,
			Path:     "/",
			SameSite: cm.someSite,
			Domain:   cm.config.DOMAIN,
		},
	}
}

func (cm *CookieManager) CreateRevokeCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "SessionId",
		Value:    "",
		HttpOnly: true,
		Secure:   cm.secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Domain:   cm.config.DOMAIN,
		MaxAge:   -1,
	}
}
