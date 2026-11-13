package study

import (
	"bytes"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/table"
)

type GradebookHandler struct {
	service GradebookService
}

func (h GradebookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	format, ok := table.ParseFormat(r.URL.Query().Get("format"))
	if !ok {
		badRequest(w)
		return
	}
	book, err := h.service.Build(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	var buf bytes.Buffer
	if err := table.Write(&buf, format, book.Header, book.Rows); err != nil {
		slog.Error("write gradebook", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
		return
	}
	w.Header().Set("Content-Type", format.ContentType())
	w.Header().Set("Content-Disposition", `attachment; filename="gradebook`+format.Extension()+`"`)
	_, _ = w.Write(buf.Bytes())
}
