package profile

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UpdateProfileHandler struct {
	service Service
}

func (h UpdateProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ProfileInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	out, err := h.service.Update(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, out)
}
