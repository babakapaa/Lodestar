package access

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type MyPermissionsHandler struct {
	service Service
}

func (h MyPermissionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	perms, err := h.service.Permissions(r.Context(), principal(r), r.URL.Query().Get("scope"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, perms)
}
