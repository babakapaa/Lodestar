package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ToolHandler struct {
	service Service
}

func (h ToolHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.Tool(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, urls)
}
