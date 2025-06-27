package handler

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// spaHandler implements http.Handler and serves files from the embedded filesystem
// with SPA routing support - it serves index.html for all non-file routes
type spaHandler struct {
	fsys fs.FS
}

// NewSPAHandler creates a new SPA handler that serves files from the given filesystem
// and falls back to index.html for client-side routing
func NewSPAHandler(fsys fs.FS) http.Handler {
	return &spaHandler{fsys: fsys}
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean the path
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	upath = path.Clean(upath)

	// Try to open the file
	file, err := h.fsys.Open(strings.TrimPrefix(upath, "/"))
	if err == nil {
		defer file.Close()
		// Check if it's a directory
		stat, err := file.Stat()
		if err == nil && !stat.IsDir() {
			// Serve the file
			http.ServeContent(w, r, upath, stat.ModTime(), file.(io.ReadSeeker))
			return
		}
	}

	// If file not found or is a directory, serve index.html
	indexFile, err := h.fsys.Open("index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer indexFile.Close()

	stat, err := indexFile.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Set proper content type for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
}
