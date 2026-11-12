package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ScheduleHandler struct {
	service ScheduleService
}

func (h ScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	from, to, err := dayRange(r)
	if err != nil {
		writeError(w, err)
		return
	}
	list, err := h.service.Mine(r.Context(), principal(r), from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
