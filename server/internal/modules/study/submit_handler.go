package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SubmitHandler struct {
	service AssignmentService
}

func (h SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SubmissionInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	sub, created, err := h.service.Submit(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	jsonhttp.Write(w, status, sub)
}
