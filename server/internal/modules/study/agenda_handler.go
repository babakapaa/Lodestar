package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type AgendaHandler struct {
	service ScheduleService
}

func (h AgendaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var date *Day
	if value := r.URL.Query().Get("date"); value != "" {
		d, err := ParseDay(value)
		if err != nil {
			writeError(w, err)
			return
		}
		date = &d
	}
	a, err := h.service.Agenda(r.Context(), principal(r), date)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, a)
}
