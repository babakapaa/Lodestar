package access

import "net/http"

type DeleteRoleHandler struct {
	service Service
}

func (h DeleteRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteRole(r.Context(), principal(r), r.PathValue("id")); err != nil {
		WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
