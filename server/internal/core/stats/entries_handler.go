package stats

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type EntriesHandler struct {
	service Service
}

func (h EntriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.PrincipalFrom(r.Context())
	list, err := h.service.Entries(r.Context(), p)
	if err != nil {
		access.WriteError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
