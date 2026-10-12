package cookie_utils

import (
	"fmt"
	"net/http"
	"net/url"
)

func EncodeCookies(cookies []*http.Cookie) []string {
	var results []string

	for _, cookie := range cookies {
		result := fmt.Sprintf("%s=%s", url.QueryEscape(cookie.Name), url.QueryEscape(cookie.Value))

		if cookie.Path != "" {
			result += "; Path=" + cookie.Path
		}

		if !cookie.Expires.IsZero() {
			result += "; Expires=" + cookie.Expires.UTC().Format(http.TimeFormat)
		}

		if cookie.MaxAge > 0 {
			result += fmt.Sprintf("; Max-Age=%d", cookie.MaxAge)
		} else if cookie.MaxAge < 0 {
			result += "; Max-Age=0"
		}

		if cookie.Domain != "" {
			result += "; Domain=" + cookie.Domain
		}

		if cookie.HttpOnly {
			result += "; HttpOnly"
		}

		if cookie.Secure {
			result += "; Secure"
		}

		switch cookie.SameSite {
		case http.SameSiteLaxMode:
			result += "; SameSite=Lax"
		case http.SameSiteStrictMode:
			result += "; SameSite=Strict"
		case http.SameSiteNoneMode:
			result += "; SameSite=None"
		}

		results = append(results, result)
	}

	return results
}
