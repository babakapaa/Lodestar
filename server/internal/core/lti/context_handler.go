package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ContextHandler struct {
	service Service
}

func (h ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.service.Context(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, v)
}
