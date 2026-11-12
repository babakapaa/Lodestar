package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UpdateCourseHandler struct {
	service Service
}

func (h UpdateCourseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in CourseInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	c, err := h.service.Update(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, c)
}
