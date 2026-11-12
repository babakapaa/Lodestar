package privacy

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Routes struct {
	Service       Service
	Exporter      Exporter
	Authenticator auth.Authenticator
}

func (r Routes) Mount(mux *http.ServeMux) {
	require := r.Authenticator.Require
	mux.Handle("GET /v1/me/export", require(ExportHandler{exporter: r.Exporter}))
	mux.Handle("GET /v1/me/deletion", require(DeletionHandler{service: r.Service}))
	mux.Handle("POST /v1/me/deletion", require(RequestDeletionHandler{service: r.Service}))
	mux.Handle("DELETE /v1/me/deletion", require(CancelDeletionHandler{service: r.Service}))
}
