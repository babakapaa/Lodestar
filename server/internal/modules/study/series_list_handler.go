package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SeriesListHandler struct {
	service ScheduleService
}

func (h SeriesListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Series(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, orEmpty(list))
}
