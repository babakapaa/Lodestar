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
	switch r.URL.Query().Get("format") {
	case "md":
		attach(w, "text/markdown; charset=utf-8", fmt.Sprintf("protocol-v%d.md", p.Version), []byte(p.Markdown()))
	case "pdf":
		doc, err := p.PDF()
		if err != nil {
			writeError(w, err)
			return
		}
		attach(w, "application/pdf", fmt.Sprintf("protocol-v%d.pdf", p.Version), doc)
	default:
		jsonhttp.Write(w, http.StatusOK, p)
	}
}

func attach(w http.ResponseWriter, contentType, name string, body []byte) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	_, _ = w.Write(body)
}
