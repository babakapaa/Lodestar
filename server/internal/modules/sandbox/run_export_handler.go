package sandbox

import (
	"bytes"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/table"
)

type ExportRunsHandler struct{ service Service }

func (h ExportRunsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	format, ok := table.ParseFormat(r.URL.Query().Get("format"))
	if !ok {
		badRequest(w)
		return
	}
	runs, err := h.service.ExportRuns(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	var buf bytes.Buffer
	if err := table.Write(&buf, format, runs.Header, runs.Rows); err != nil {
		slog.Error("write sandbox runs", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
		return
	}
	w.Header().Set("Content-Type", format.ContentType())
	w.Header().Set("Content-Disposition", `attachment; filename="gpu-runs`+format.Extension()+`"`)
	_, _ = w.Write(buf.Bytes())
}
