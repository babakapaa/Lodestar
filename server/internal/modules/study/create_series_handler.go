package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CreateSeriesHandler struct {
	service ScheduleService
}

func (h CreateSeriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SeriesInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		writeError(w, ErrInvalidSeries)
		return
	}
	s, err := h.service.Create(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusCreated, s)
}
