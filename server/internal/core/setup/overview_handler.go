package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type OverviewHandler struct {
	service Service
}

func (h OverviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o, err := h.service.Overview(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, o)
}
