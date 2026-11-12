package study

import "net/http"

type DeleteAssignmentHandler struct {
	service AssignmentService
}

func (h DeleteAssignmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
