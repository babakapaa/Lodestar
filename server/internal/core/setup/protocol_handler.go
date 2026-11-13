package setup

import (
	"fmt"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ProtocolHandler struct {
	service Service
}

func (h ProtocolHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, err := h.service.Protocol(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	if r.URL.Query().Get("format") != "md" {
		jsonhttp.Write(w, http.StatusOK, p)
		return
	}
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="protocol-v%d.md"`, p.Version))
	_, _ = w.Write([]byte(p.Markdown()))
}
