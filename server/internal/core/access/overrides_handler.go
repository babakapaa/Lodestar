package access

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type OverridesHandler struct {
	service Service
}

func (h OverridesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Overrides(r.Context(), principal(r), r.URL.Query().Get("scope"))
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []Override{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
