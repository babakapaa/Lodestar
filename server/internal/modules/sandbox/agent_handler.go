package sandbox

import (
	"net/http"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type AgentHandler struct{ service Service }

func (h AgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok {
		writeError(w, ErrBadNodeToken)
		return
	}
	var report AgentReport
	if err := jsonhttp.Read(w, r, &report); err != nil {
		badRequest(w)
		return
	}
	plan, err := h.service.Sync(r.Context(), raw, report)
	respond(w, http.StatusOK, plan, err)
}
