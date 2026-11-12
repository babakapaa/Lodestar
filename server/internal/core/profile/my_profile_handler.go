package profile

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type MyProfileHandler struct {
	service Service
}

func (h MyProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.Me(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, out)
}
