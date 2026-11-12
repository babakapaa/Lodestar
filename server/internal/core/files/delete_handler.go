package files

import "net/http"

type DeleteHandler struct {
	service Service
}

func (h DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
