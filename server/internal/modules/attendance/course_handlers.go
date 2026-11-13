package attendance

import (
	"bytes"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/table"
)

type CourseHandler struct{ service Service }

func (h CourseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	c, err := h.service.Course(r.Context(), principal(r), r.PathValue("id"), q.Get("from"), q.Get("to"))
	respond(w, c, err)
}

type ExportHandler struct{ service Service }

func (h ExportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	format, ok := table.ParseFormat(q.Get("format"))
	if !ok {
		badRequest(w)
		return
	}
	t, err := h.service.ExportCourse(r.Context(), principal(r), r.PathValue("id"), q.Get("from"), q.Get("to"))
	if err != nil {
		writeError(w, err)
		return
	}
	var buf bytes.Buffer
	if err := table.Write(&buf, format, t.Header, t.Rows); err != nil {
		slog.Error("write attendance", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
		return
	}
	w.Header().Set("Content-Type", format.ContentType())
	w.Header().Set("Content-Disposition", `attachment; filename="attendance`+format.Extension()+`"`)
	_, _ = w.Write(buf.Bytes())
}
