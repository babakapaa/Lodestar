package profile

import "net/http"

type ResetPreferenceHandler struct {
	service Service
}

func (h ResetPreferenceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.ResetPreference(r.Context(), principal(r), r.PathValue("key")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
