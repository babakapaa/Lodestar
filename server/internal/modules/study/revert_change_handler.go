package study

import "net/http"

type RevertChangeHandler struct {
	service ScheduleService
}

func (h RevertChangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	date, err := ParseDay(r.PathValue("date"))
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.service.RevertChange(r.Context(), principal(r), r.PathValue("id"), date); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
