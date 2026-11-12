package structure

import "net/http"

type MemberChangeHandler struct {
	service Service
	add     bool
}

func (h MemberChangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	change := h.service.RemoveMember
	if h.add {
		change = h.service.AddMember
	}
	if err := change(r.Context(), principal(r), r.PathValue("id"), r.PathValue("user")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
