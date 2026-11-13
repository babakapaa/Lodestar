package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SourceHandler struct {
	service Service
}

func (h SourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SourceInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	v, err := h.service.SetSource(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, v)
}
