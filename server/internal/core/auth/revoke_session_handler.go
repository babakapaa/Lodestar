package auth

import "net/http"

type RevokeSessionHandler struct {
	sessions SessionStore
}

func (h RevokeSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, _ := PrincipalFrom(r.Context())
	if err := h.sessions.Revoke(r.Context(), p, r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
