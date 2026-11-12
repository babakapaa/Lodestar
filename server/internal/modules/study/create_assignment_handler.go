package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CreateAssignmentHandler struct {
	service AssignmentService
}

func (h CreateAssignmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in AssignmentInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	a, err := h.service.Create(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusCreated, a)
}
