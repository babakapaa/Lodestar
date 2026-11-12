package delivery

import "net/http"

type RemoveHandler struct {
	service Service
}

func (h RemoveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Remove(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
