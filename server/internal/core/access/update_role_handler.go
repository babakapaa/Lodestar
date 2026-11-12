package access

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UpdateRoleHandler struct {
	service Service
}

func (h UpdateRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in RoleInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	role, err := h.service.UpdateRole(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, role)
}
