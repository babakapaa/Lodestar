package auth

import "net/http"

type LogoutHandler struct {
	sessions SessionStore
}

func (h LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, _ := PrincipalFrom(r.Context())
	if err := h.sessions.Revoke(r.Context(), p, p.SessionID); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
