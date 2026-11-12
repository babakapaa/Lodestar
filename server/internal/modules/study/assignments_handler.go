package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type AssignmentsHandler struct {
	service AssignmentService
}

func (h AssignmentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
