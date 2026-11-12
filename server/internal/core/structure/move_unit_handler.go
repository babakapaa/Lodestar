package structure

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type MoveUnitHandler struct {
	service Service
}

func (h MoveUnitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in MoveInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	if err := h.service.Move(r.Context(), principal(r), r.PathValue("id"), in.ParentID); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
