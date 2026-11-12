package delivery

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type DevicesHandler struct {
	service Service
}

func (h DevicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Devices(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []Device{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
