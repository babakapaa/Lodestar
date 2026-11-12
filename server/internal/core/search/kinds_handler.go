package search

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type KindsHandler struct {
	service Service
}

func (h KindsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	kinds, err := h.service.Kinds(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, kinds)
}
