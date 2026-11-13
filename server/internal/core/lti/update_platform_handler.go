package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UpdatePlatformHandler struct {
	service Service
}

func (h UpdatePlatformHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in PlatformInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	p, err := h.service.UpdatePlatform(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, p)
}
