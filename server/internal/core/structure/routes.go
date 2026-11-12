package structure

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
	mux.Handle("GET /v1/units", require(UnitsHandler{service: s}))
	mux.Handle("POST /v1/units", require(CreateUnitHandler{service: s}))
	mux.Handle("POST /v1/units/import", require(ImportHandler{run: s.ImportUnits}))
	mux.Handle("POST /v1/units/import-members", require(ImportHandler{run: s.ImportMembers}))
	mux.Handle("GET /v1/units/{id}", require(UnitHandler{service: s}))
	mux.Handle("PUT /v1/units/{id}", require(UpdateUnitHandler{service: s}))
	mux.Handle("POST /v1/units/{id}/move", require(MoveUnitHandler{service: s}))
	mux.Handle("POST /v1/units/{id}/archive", require(ArchiveHandler{service: s, archived: true}))
	mux.Handle("POST /v1/units/{id}/restore", require(ArchiveHandler{service: s}))
	mux.Handle("GET /v1/units/{id}/members", require(MembersHandler{service: s}))
	mux.Handle("PUT /v1/units/{id}/members/{user}", require(MemberChangeHandler{service: s, add: true}))
	mux.Handle("DELETE /v1/units/{id}/members/{user}", require(MemberChangeHandler{service: s}))
	mux.Handle("GET /v1/units/{id}/history", require(HistoryHandler{service: s}))
	mux.Handle("GET /v1/users/{id}/units", require(UserUnitsHandler{service: s}))
}
