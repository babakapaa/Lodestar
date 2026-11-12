package delivery

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
	mux.Handle("GET /v1/me/devices", require(DevicesHandler{service: s}))
	mux.Handle("POST /v1/me/devices", require(RegisterHandler{service: s}))
	mux.Handle("DELETE /v1/me/devices/{id}", require(RemoveHandler{service: s}))
}
