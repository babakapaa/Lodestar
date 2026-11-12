package app

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/health"
	"github.com/babakapaa/lodestar/server/internal/core/meta"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func NewRouter(version string, registry *module.Registry, db health.Pinger, gate module.Gate, kernel module.Kernel, core ...Routes) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /healthz", health.Handler{})
	mux.Handle("GET /readyz", health.NewReadyHandler(db))
	mux.Handle("GET /v1/meta", meta.NewHandler(version, registry))
	for _, routes := range core {
		routes.Mount(mux)
	}
	registry.Mount(mux, gate, kernel)
	mux.HandleFunc("/", apierror.NotFound)
	return mux
}
