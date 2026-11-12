package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UpdateSeriesHandler struct {
	service ScheduleService
}

func (h UpdateSeriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SeriesInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		writeError(w, ErrInvalidSeries)
		return
	}
	s, err := h.service.Update(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, s)
}
