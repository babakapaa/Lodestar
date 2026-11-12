package settings

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
	mux.Handle("GET /v1/settings/schema", require(SchemaHandler{service: s}))
	mux.Handle("GET /v1/settings", require(ValuesHandler{service: s}))
	mux.Handle("PUT /v1/settings/{key}", require(SetHandler{service: s}))
	mux.Handle("DELETE /v1/settings/{key}", require(ResetHandler{service: s}))
	mux.Handle("GET /v1/modules", require(ModulesHandler{service: s}))
	mux.Handle("PUT /v1/modules/{name}", require(SetModuleHandler{service: s}))
}
