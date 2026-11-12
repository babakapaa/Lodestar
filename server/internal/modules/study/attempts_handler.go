package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type AttemptsHandler struct {
	service AssignmentService
}

func (h AttemptsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Attempts(r.Context(), principal(r), r.PathValue("id"), r.PathValue("student"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
