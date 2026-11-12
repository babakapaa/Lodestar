package privacy

import "net/http"

type CancelDeletionHandler struct {
	service Service
}

func (h CancelDeletionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.CancelDeletion(r.Context(), principal(r)); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
