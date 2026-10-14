package middleware

import (
	"net/http"
	"net/url"
)

type Middleware = func(http.Handler) http.Handler

type Server[R Route] interface {
	FindPath(method string, u *url.URL) (r R, _ bool)
}

type Route interface {
	Name() string
	OperationID() string
	PathPattern() string
}

// RouteFinder finds Route by given URL.
type RouteFinder func(method string, u *url.URL) (Route, bool)

// MakeRouteFinder creates RouteFinder from given server.
func MakeRouteFinder[R Route, S Server[R]](server S) RouteFinder {
	return func(method string, u *url.URL) (Route, bool) {
		return server.FindPath(method, u)
	}
}

// urlWithoutQuery は認証情報・クエリ・フラグメントを除いたURLを返す。
func urlWithoutQuery(u *url.URL) string {
	stripped := *u
	stripped.User = nil
	stripped.RawQuery = ""
	stripped.ForceQuery = false
	stripped.Fragment = ""
	stripped.RawFragment = ""
	return stripped.String()
}
