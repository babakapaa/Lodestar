package profile

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SetPreferenceHandler struct {
	service Service
}

func (h SetPreferenceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ValueInput
	if err := jsonhttp.Read(w, r, &in); err != nil || len(in.Value) == 0 {
		badRequest(w)
		return
	}
	if err := h.service.SetPreference(r.Context(), principal(r), r.PathValue("key"), in.Value); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
