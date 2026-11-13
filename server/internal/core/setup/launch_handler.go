package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type LaunchHandler struct {
	service Service
}

func (h LaunchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in LaunchInput
	if r.ContentLength != 0 {
		if err := jsonhttp.Read(w, r, &in); err != nil {
			badRequest(w)
			return
		}
	}
	o, err := h.service.Launch(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, o)
}
