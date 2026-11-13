package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CreatePlatformHandler struct {
	service Service
}

func (h CreatePlatformHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in PlatformInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	p, err := h.service.CreatePlatform(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusCreated, p)
}
