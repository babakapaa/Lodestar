package access

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
	mux.Handle("GET /v1/permissions", require(CatalogHandler{service: s}))
	mux.Handle("GET /v1/me/permissions", require(MyPermissionsHandler{service: s}))
	mux.Handle("GET /v1/roles", require(RolesHandler{service: s}))
	mux.Handle("POST /v1/roles", require(CreateRoleHandler{service: s}))
	mux.Handle("PUT /v1/roles/{id}", require(UpdateRoleHandler{service: s}))
	mux.Handle("DELETE /v1/roles/{id}", require(DeleteRoleHandler{service: s}))
	mux.Handle("POST /v1/roles/{id}/move", require(MoveRoleHandler{service: s}))
	mux.Handle("GET /v1/roles/{id}/preview", require(PreviewHandler{service: s}))
	mux.Handle("GET /v1/users/{id}/roles", require(UserRolesHandler{service: s}))
	mux.Handle("POST /v1/role-assignments", require(AssignHandler{service: s}))
	mux.Handle("DELETE /v1/role-assignments/{id}", require(UnassignHandler{service: s}))
	mux.Handle("GET /v1/overrides", require(OverridesHandler{service: s}))
	mux.Handle("PUT /v1/overrides/{role}", require(SetOverrideHandler{service: s}))
}
