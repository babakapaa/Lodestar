package lti

import "net/http"

type DeletePlatformHandler struct {
	service Service
}

func (h DeletePlatformHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeletePlatform(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
