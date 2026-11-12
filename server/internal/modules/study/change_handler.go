package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ChangeHandler struct {
	service ScheduleService
}

func (h ChangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	date, err := ParseDay(r.PathValue("date"))
	if err != nil {
		writeError(w, err)
		return
	}
	var in ChangeInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		writeError(w, ErrInvalidChange)
		return
	}
	ch, err := h.service.PutChange(r.Context(), principal(r), r.PathValue("id"), date, in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, ch)
}
