package access

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SetOverrideHandler struct {
	service Service
}

func (h SetOverrideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in OverrideInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	o := Override{RoleID: r.PathValue("role"), Allow: in.Allow, Deny: in.Deny}
	if err := h.service.SetOverride(r.Context(), principal(r), r.URL.Query().Get("scope"), o); err != nil {
		WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
