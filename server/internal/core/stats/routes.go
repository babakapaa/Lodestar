package stats

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Routes struct {
	Service       Service
	Authenticator auth.Authenticator
}

func (r Routes) Mount(mux *http.ServeMux) {
	mux.Handle("GET /v1/stats/units", r.Authenticator.Require(UnitsHandler{service: r.Service}))
}
