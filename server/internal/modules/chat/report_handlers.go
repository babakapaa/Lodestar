package chat

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ReportsHandler struct{ service ReportService }

func (h ReportsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Open(r.Context(), principal(r))
	writeOK(w, list, err)
}

type ReportHandler struct{ service ReportService }

func (h ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ReportInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	report, err := h.service.Create(r.Context(), principal(r), in)
	writeCreated(w, report, true, err)
}

type ResolveHandler struct{ service ReportService }

func (h ResolveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.Resolve(r.Context(), principal(r), r.PathValue("id"))
	writeOK(w, report, err)
}
