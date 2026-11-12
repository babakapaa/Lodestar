package structure

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CreateUnitHandler struct {
	service Service
}

func (h CreateUnitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in UnitInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	u, err := h.service.Create(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusCreated, u)
}
