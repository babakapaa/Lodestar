package search

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
	mux.Handle("GET /v1/search", require(SearchHandler{service: r.Service}))
	mux.Handle("GET /v1/search/kinds", require(KindsHandler{service: r.Service}))
}
