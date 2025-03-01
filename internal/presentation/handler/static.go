package handler

import (
	"io/fs"
	"net/http"

	"github.com/neko-dream/server/static"
)

func NewStaticHandler() http.Handler {
	fsys, err := fs.Sub(static.Static, ".")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
}
