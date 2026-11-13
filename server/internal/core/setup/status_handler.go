package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type StatusHandler struct {
	service Service
}

func (h StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	view, err := h.service.Status(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, view)
}
