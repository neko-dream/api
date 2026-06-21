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

		if cookie.SameSite != http.SameSiteDefaultMode {
			result += fmt.Sprintf("; SameSite=%v", "Lax")
		}

		results = append(results, result)
	}

	return results
}
