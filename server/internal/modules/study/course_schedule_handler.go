package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CourseScheduleHandler struct {
	service ScheduleService
}

func (h CourseScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	from, to, err := dayRange(r)
	if err != nil {
		writeError(w, err)
		return
	}
	list, err := h.service.OfCourse(r.Context(), principal(r), r.PathValue("id"), from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
