package structure

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UpdateUnitHandler struct {
	service Service
}

func (h UpdateUnitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in UnitInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	u, err := h.service.Update(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, u)
}
