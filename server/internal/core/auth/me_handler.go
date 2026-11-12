package auth

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type MeHandler struct {
	users identity.Store
}

func (h MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, _ := PrincipalFrom(r.Context())
	user, err := h.users.ByID(r.Context(), p.TenantID, p.UserID)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, user)
}
