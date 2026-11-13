package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type PlatformsHandler struct {
	service Service
}

func (h PlatformsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Platforms(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
