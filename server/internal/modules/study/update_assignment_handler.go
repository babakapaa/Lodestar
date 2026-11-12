package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UpdateAssignmentHandler struct {
	service AssignmentService
}

func (h UpdateAssignmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in AssignmentInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	a, err := h.service.Update(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, a)
}
