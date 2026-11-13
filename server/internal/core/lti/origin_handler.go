package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type OriginHandler struct {
	service Service
}

func (h OriginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.ByOrigin(r.Context(), principal(r), r.PathValue("kind"), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
