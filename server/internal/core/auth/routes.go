package auth

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

type Routes struct {
	Login         LoginService
	Sessions      SessionStore
	Users         identity.Store
	Authenticator Authenticator
}

func (r Routes) Mount(mux *http.ServeMux) {
	require := r.Authenticator.Require
	mux.Handle("GET /v1/auth/{tenant}/providers", ProvidersHandler{login: r.Login})
	mux.Handle("GET /v1/auth/{tenant}/{provider}/login", LoginHandler{login: r.Login})
	mux.Handle("GET /v1/auth/{tenant}/{provider}/callback", CallbackHandler{login: r.Login})
	mux.Handle("POST /v1/auth/token", TokenHandler{login: r.Login, sessions: r.Sessions})
	mux.Handle("POST /v1/auth/logout", require(LogoutHandler{sessions: r.Sessions}))
	mux.Handle("GET /v1/me", require(MeHandler{users: r.Users}))
	mux.Handle("GET /v1/sessions", require(SessionsHandler{sessions: r.Sessions}))
	mux.Handle("DELETE /v1/sessions/{id}", require(RevokeSessionHandler{sessions: r.Sessions}))
}
