package study

import "net/http"

type RosterChangeHandler struct {
	service Service
	list    List
	add     bool
}

func (h RosterChangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.service.Change(r.Context(), principal(r), h.list, r.PathValue("id"), r.PathValue("member"), h.add)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
