package audit

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Routes struct {
	Service       Service
	Authenticator auth.Authenticator
}

func (r Routes) Mount(mux *http.ServeMux) {
	mux.Handle("GET /v1/audit", r.Authenticator.Require(ListHandler{service: r.Service}))
}
