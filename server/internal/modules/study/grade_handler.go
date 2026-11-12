package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type GradeHandler struct {
	service AssignmentService
}

func (h GradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in GradeInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	g, err := h.service.Grade(r.Context(), principal(r), r.PathValue("id"), r.PathValue("student"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, g)
}
