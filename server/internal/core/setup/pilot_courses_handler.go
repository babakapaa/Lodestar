package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type PilotCoursesHandler struct {
	service Service
}

func (h PilotCoursesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.PilotCourses(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
