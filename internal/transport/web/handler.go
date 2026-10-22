package web

import (
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
)

// Handler serves the embedded single-page application. Any path that does not
// name a real file falls back to index.html so client-side routes survive a
// page refresh.
type Handler struct {
	dist     fs.FS
	index    []byte
	indexErr error
}

// NewHandler wraps the embedded frontend tree.
func NewHandler(dist fs.FS) *Handler {
	h := &Handler{dist: dist}
	h.index, h.indexErr = fs.ReadFile(dist, "index.html")
	return h
}

// ServeHTTP serves static files, falling back to index.html for app routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if name != "" && name != "." {
		if content, err := fs.ReadFile(h.dist, name); err == nil {
			if strings.HasPrefix(name, "assets/") {
				// Vite fingerprints asset names, so they never change.
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			writeFile(w, r, name, content)
			return
		}
	}

	if h.indexErr != nil {
		http.Error(w, "前端资源未打包进当前二进制，请先执行 make build 或 make frontend-build", http.StatusNotFound)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	writeFile(w, r, "index.html", h.index)
}

func writeFile(w http.ResponseWriter, r *http.Request, name string, content []byte) {
	if contentType := mime.TypeByExtension(path.Ext(name)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}
	_, _ = w.Write(content)
}
