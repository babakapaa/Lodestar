package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type StudentsHandler struct {
	service Service
}

func (h StudentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Students(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
