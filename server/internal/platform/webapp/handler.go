package webapp

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var apiPaths = []string{"/v1/", "/healthz", "/readyz"}

type Handler struct {
	root  string
	files http.Handler
	api   http.Handler
}

func New(root string, api http.Handler) Handler {
	return Handler{root: root, files: http.FileServer(http.Dir(root)), api: api}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !servesPage(r) {
		h.api.ServeHTTP(w, r)
		return
	}
	name := path.Clean(r.URL.Path)
	if name != "/" && h.isFile(name) {
		h.files.ServeHTTP(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, filepath.Join(h.root, "index.html"))
}

func (h Handler) isFile(name string) bool {
	info, err := os.Stat(filepath.Join(h.root, filepath.FromSlash(name)))
	return err == nil && !info.IsDir()
}

func servesPage(r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	for _, prefix := range apiPaths {
		if strings.HasPrefix(r.URL.Path, prefix) {
			return false
		}
	}
	return true
}
