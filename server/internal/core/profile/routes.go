package profile

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
	mux.Handle("GET /v1/me/profile", require(MyProfileHandler{service: s}))
	mux.Handle("PUT /v1/me/profile", require(UpdateProfileHandler{service: s}))
	mux.Handle("GET /v1/users/{id}/profile", require(UserProfileHandler{service: s}))
	mux.Handle("GET /v1/preferences/schema", require(PreferenceSchemaHandler{service: s}))
	mux.Handle("GET /v1/me/preferences", require(PreferencesHandler{service: s}))
	mux.Handle("PUT /v1/me/preferences/{key}", require(SetPreferenceHandler{service: s}))
	mux.Handle("DELETE /v1/me/preferences/{key}", require(ResetPreferenceHandler{service: s}))
}
