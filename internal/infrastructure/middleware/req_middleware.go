package middleware

import (
	"log"
	"net/http"

	http_utils "github.com/neko-dream/server/pkg/http"
)

func ReqMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(http_utils.WithHTTPResReqContext(r.Context(), r, w))
		log.Println("Request URL: ", r.URL)
		log.Println("COOKIE HEADER: ", r.Header.Get("Cookie"))
		next.ServeHTTP(w, r)
	})
}
