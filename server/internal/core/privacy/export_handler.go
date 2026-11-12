package privacy

import (
	"log/slog"
	"net/http"
	"time"
)

type ExportHandler struct {
	exporter Exporter
}

func (h ExportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	archive, err := h.exporter.Collect(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="lodestar-export-`+time.Now().UTC().Format("2006-01-02")+`.zip"`)
	if err := archive.Write(r.Context(), w); err != nil {
		slog.Error("export personal data", "error", err)
		panic(http.ErrAbortHandler)
	}
}
