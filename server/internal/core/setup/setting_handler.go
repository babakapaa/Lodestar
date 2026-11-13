package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SettingHandler struct {
	service Service
}

func (h SettingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SettingInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	err := h.service.SetSetting(r.Context(), principal(r), r.PathValue("section"), r.PathValue("key"), in.Value)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
