package auth

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SessionsHandler struct {
	sessions SessionStore
}

func (h SessionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, _ := PrincipalFrom(r.Context())
	list, err := h.sessions.List(r.Context(), p)
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []Session{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
