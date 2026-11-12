package settings

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ModulesHandler struct {
	service Service
}

func (h ModulesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	states, err := h.service.Modules(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, states)
}
