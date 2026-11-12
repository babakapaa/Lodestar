package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CourseHandler struct {
	service Service
}

func (h CourseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Get(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, d)
}
