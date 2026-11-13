package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UnlinkHandler struct {
	service Service
}

func (h UnlinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.service.Unlink(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, v)
}
