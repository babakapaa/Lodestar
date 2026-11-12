package structure

import (
	"context"
	"io"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const maxImport = 16 << 20

type ImportHandler struct {
	run func(ctx context.Context, p auth.Principal, r io.Reader) (ImportReport, error)
}

func (h ImportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	report, err := h.run(r.Context(), principal(r), http.MaxBytesReader(w, r.Body, maxImport))
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
