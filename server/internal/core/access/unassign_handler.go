package access

import "net/http"

type UnassignHandler struct {
	service Service
}

func (h UnassignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Unassign(r.Context(), principal(r), r.PathValue("id")); err != nil {
		WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
