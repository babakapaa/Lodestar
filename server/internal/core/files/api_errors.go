package files

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeInvalidRequest = "invalid_request"
	CodeTooLarge       = "file_too_large"
	CodeNotFound       = "file_not_found"
	CodeNotUploaded    = "file_not_uploaded"
	CodeSizeMismatch   = "file_size_mismatch"
	CodeUnavailable    = "files_unavailable"
	CodeInUse          = "file_in_use"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidFile):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
	case errors.Is(err, ErrTooLarge):
		apierror.Write(w, http.StatusRequestEntityTooLarge, CodeTooLarge)
	case errors.Is(err, ErrNotFound):
		apierror.Write(w, http.StatusNotFound, CodeNotFound)
	case errors.Is(err, ErrNotUploaded):
		apierror.Write(w, http.StatusConflict, CodeNotUploaded)
	case errors.Is(err, ErrSizeMismatch):
		apierror.Write(w, http.StatusConflict, CodeSizeMismatch)
	case errors.Is(err, ErrInUse):
		apierror.Write(w, http.StatusConflict, CodeInUse)
	case errors.Is(err, ErrUnavailable):
		apierror.Write(w, http.StatusServiceUnavailable, CodeUnavailable)
	default:
		slog.Error("file request failed", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
	}
}

func badRequest(w http.ResponseWriter) {
	apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
}
