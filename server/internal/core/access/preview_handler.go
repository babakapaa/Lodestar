package access

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type PreviewHandler struct {
	service Service
}

func (h PreviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	perms, err := h.service.Preview(r.Context(), principal(r), r.PathValue("id"), r.URL.Query().Get("scope"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, perms)
}
