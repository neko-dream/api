package handler

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/neko-dream/server/static"
)

func NewStaticHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/static/")
		// Check if it's a Web Push test file
		if path == "test-webpush.html" || path == "firebase-messaging-sw.js" {
			fsys, err := fs.Sub(static.WebPushTest, ".")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			newReq := r.Clone(r.Context())
			newReq.URL.Path = "/" + path
			http.FileServer(http.FS(fsys)).ServeHTTP(w, newReq)
			return
		}
		// Default to OAS files
		fsys, err := fs.Sub(static.Oas, ".")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Create new request with modified path
		newReq := r.Clone(r.Context())
		newReq.URL.Path = "/" + path
		http.FileServer(http.FS(fsys)).ServeHTTP(w, newReq)
	})
}

func NewAdminUIHandler() http.Handler {
	fsys, err := fs.Sub(static.AdminUI, "admin-ui")
	if err != nil {
		panic(err)
	}
	return NewSPAHandler(fsys)
}

func NewAdminUIAssetsHandler() http.Handler {
	fsys, err := fs.Sub(static.AdminUI, "admin-ui/assets")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
}
