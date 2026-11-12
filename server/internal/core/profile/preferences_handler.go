package profile

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type PreferencesHandler struct {
	service Service
}

func (h PreferencesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values, err := h.service.Preferences(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, values)
}
