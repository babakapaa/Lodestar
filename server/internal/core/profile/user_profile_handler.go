package profile

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UserProfileHandler struct {
	service Service
}

func (h UserProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.Of(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, out)
}
