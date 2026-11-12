package files

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Routes struct {
	Service       Service
	Authenticator auth.Authenticator
}

func (r Routes) Mount(mux *http.ServeMux) {
	require := r.Authenticator.Require
	s := r.Service
	mux.Handle("GET /v1/files", require(ListHandler{service: s}))
	mux.Handle("POST /v1/files", require(CreateHandler{service: s}))
	mux.Handle("GET /v1/files/{id}", require(FileHandler{service: s}))
	mux.Handle("DELETE /v1/files/{id}", require(DeleteHandler{service: s}))
	mux.Handle("POST /v1/files/{id}/complete", require(CompleteHandler{service: s}))
	mux.Handle("GET /v1/files/{id}/download", require(DownloadHandler{service: s}))
}
