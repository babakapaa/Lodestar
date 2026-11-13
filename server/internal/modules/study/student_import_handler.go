package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const maxImport = 16 << 20

type StudentImportHandler struct {
	service StudentImportService
}

func (h StudentImportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.Import(r.Context(), principal(r), r.PathValue("id"), http.MaxBytesReader(w, r.Body, maxImport))
	if err != nil {
		writeError(w, err)
		return
	}
	status := http.StatusOK
	if len(report.Errors) > 0 {
		status = http.StatusUnprocessableEntity
	}
	jsonhttp.Write(w, status, report)
}
