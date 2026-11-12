package notify

import "net/http"

type ReadHandler struct {
	service Service
}

func (h ReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.MarkRead(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
