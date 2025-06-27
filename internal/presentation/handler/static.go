package handler

import (
	"io/fs"
	"net/http"

	"github.com/neko-dream/server/static"
)

func NewStaticHandler() http.Handler {
	fsys, err := fs.Sub(static.Oas, ".")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
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
