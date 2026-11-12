package profile

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	CodeUserNotFound   = "user_not_found"
	CodeInvalidProfile = "invalid_profile"
	CodeUnknownKey     = "unknown_key"
	CodeInvalidValue   = "invalid_value"
	CodeInvalidRequest = "invalid_request"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		apierror.Write(w, http.StatusNotFound, CodeUserNotFound)
	case errors.Is(err, ErrInvalidProfile):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidProfile)
	case errors.Is(err, schema.ErrUnknownKey):
		apierror.Write(w, http.StatusNotFound, CodeUnknownKey)
	case errors.Is(err, schema.ErrInvalidValue):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidValue)
	default:
		slog.Error("profile request failed", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
	}
}

func badRequest(w http.ResponseWriter) {
	apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
}
