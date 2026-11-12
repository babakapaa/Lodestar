package notify

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeUnknownType    = "unknown_notification_type"
	CodeRequiredType   = "notification_required"
	CodeNotFound       = "notification_not_found"
	CodeInvalidRequest = "invalid_request"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnknownType):
		apierror.Write(w, http.StatusNotFound, CodeUnknownType)
	case errors.Is(err, ErrRequiredType):
		apierror.Write(w, http.StatusBadRequest, CodeRequiredType)
	case errors.Is(err, ErrNotFound):
		apierror.Write(w, http.StatusNotFound, CodeNotFound)
	default:
		slog.Error("notification request failed", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
	}
}

func badRequest(w http.ResponseWriter) {
	apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
}
