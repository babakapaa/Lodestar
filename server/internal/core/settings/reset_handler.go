package settings

import "net/http"

type ResetHandler struct {
	service Service
}

func (h ResetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Reset(r.Context(), principal(r), r.PathValue("key")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
