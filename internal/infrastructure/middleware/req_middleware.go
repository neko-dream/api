package middleware

import (
	"net/http"

	http_utils "github.com/neko-dream/server/pkg/http"
)

func ReqMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(http_utils.WithHTTPResReqContext(r.Context(), r, w))
		next.ServeHTTP(w, r)
	})
}
