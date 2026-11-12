package notify

import "net/http"

type ReadAllHandler struct {
	service Service
}

func (h ReadAllHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.MarkAllRead(r.Context(), principal(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
