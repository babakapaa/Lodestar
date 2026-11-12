package live

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Routes struct {
	Hub           *Hub
	Tenants       database.TenantTx
	Authenticator auth.Authenticator
	Origins       []string
}

func (r Routes) Mount(mux *http.ServeMux) {
	backlog := NewBacklog(r.Tenants)
	mux.Handle("GET /v1/live", r.Authenticator.RequireStream(StreamHandler{hub: r.Hub, backlog: backlog, origins: r.Origins}))
	mux.Handle("GET /v1/events", r.Authenticator.Require(EventsHandler{backlog: backlog}))
}
