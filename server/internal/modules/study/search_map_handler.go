package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SearchMapHandler struct {
	service AssignmentService
}

func (h SearchMapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.SearchMap(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, d)
}
