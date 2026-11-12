package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CreateCourseHandler struct {
	service Service
}

func (h CreateCourseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in CourseInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	c, err := h.service.Create(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusCreated, c)
}
