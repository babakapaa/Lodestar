package settings

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ValuesHandler struct {
	service Service
}

func (h ValuesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values, err := h.service.Values(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, values)
}
