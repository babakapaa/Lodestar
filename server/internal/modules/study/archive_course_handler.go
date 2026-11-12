package study

import "net/http"

type ArchiveCourseHandler struct {
	service  Service
	archived bool
}

func (h ArchiveCourseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.SetArchived(r.Context(), principal(r), r.PathValue("id"), h.archived); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
