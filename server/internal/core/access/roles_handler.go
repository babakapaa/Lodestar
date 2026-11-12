package access

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type RolesHandler struct {
	service Service
}

func (h RolesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.Roles(r.Context(), principal(r))
	if err != nil {
		WriteError(w, err)
		return
	}
	if roles == nil {
		roles = []Role{}
	}
	jsonhttp.Write(w, http.StatusOK, roles)
}
