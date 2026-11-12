package app

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/health"
	"github.com/babakapaa/lodestar/server/internal/core/meta"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func NewRouter(version string, registry *module.Registry) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /healthz", health.Handler{})
	mux.Handle("GET /v1/meta", meta.NewHandler(version, registry))
	registry.Mount(mux)
	mux.HandleFunc("/", apierror.NotFound)
	return mux
}
