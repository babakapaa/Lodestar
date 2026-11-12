package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type MyCoursesHandler struct {
	service Service
}

func (h MyCoursesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Mine(r.Context(), principal(r), r.URL.Query().Get("archived") == "true")
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []MyCourse{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
