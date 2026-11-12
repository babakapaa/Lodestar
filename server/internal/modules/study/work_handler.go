package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type WorkHandler struct {
	service AssignmentService
}

func (h WorkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Work(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
