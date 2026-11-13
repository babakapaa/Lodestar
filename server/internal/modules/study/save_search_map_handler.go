package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SaveSearchMapHandler struct {
	service AssignmentService
}

func (h SaveSearchMapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SearchMap
	if err := jsonhttp.ReadLimit(w, r, &in, jsonhttp.LargeBody); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.SaveSearchMap(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, d)
}
