package navigator

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type HelpHandler struct{ service Service }

func (h HelpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	out, err := h.service.Help(r.Context(), principal(r), q.Get("scopeKind"), q.Get("scopeId"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, out)
}
