package access

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CreateRoleHandler struct {
	service Service
}

func (h CreateRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in RoleInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	role, err := h.service.CreateRole(r.Context(), principal(r), in)
	if err != nil {
		WriteError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusCreated, role)
}
