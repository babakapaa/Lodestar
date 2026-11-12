package access

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type MoveRoleHandler struct {
	service Service
}

func (h MoveRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in MoveInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	if err := h.service.MoveRole(r.Context(), principal(r), r.PathValue("id"), in.Position); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
