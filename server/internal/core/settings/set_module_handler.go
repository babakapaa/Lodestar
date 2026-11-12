package settings

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SetModuleHandler struct {
	service Service
}

func (h SetModuleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ModuleInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	if err := h.service.SetModule(r.Context(), principal(r), r.PathValue("name"), in.Enabled); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
