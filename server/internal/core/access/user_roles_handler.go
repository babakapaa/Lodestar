package access

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UserRolesHandler struct {
	service Service
}

func (h UserRolesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.UserRoles(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		WriteError(w, err)
		return
	}
	if list == nil {
		list = []Assignment{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
