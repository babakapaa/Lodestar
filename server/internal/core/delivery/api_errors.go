package delivery

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeInvalidRequest = "invalid_request"
	CodeDeviceNotFound = "device_not_found"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidDevice):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
	case errors.Is(err, ErrDeviceNotFound):
		apierror.Write(w, http.StatusNotFound, CodeDeviceNotFound)
	default:
		slog.Error("device request failed", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
	}
}
