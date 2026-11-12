package notify

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
	mux.Handle("GET /v1/notifications", require(ListHandler{service: s}))
	mux.Handle("GET /v1/notifications/unread", require(UnreadHandler{service: s}))
	mux.Handle("POST /v1/notifications/read-all", require(ReadAllHandler{service: s}))
	mux.Handle("POST /v1/notifications/{id}/read", require(ReadHandler{service: s}))
	mux.Handle("GET /v1/me/notification-channels", require(ChannelsHandler{service: s}))
	mux.Handle("PUT /v1/me/notification-channels/{type}", require(SetChannelsHandler{service: s}))
}
