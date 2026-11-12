package meta

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type Handler struct {
	version  string
	registry *module.Registry
}

func NewHandler(version string, registry *module.Registry) Handler {
	return Handler{version: version, registry: registry}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	jsonhttp.Write(w, http.StatusOK, Response{
		Version: h.version,
		Modules: h.registry.Manifests(),
	})
}
