package privacy

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeSelfDeletionOff = "deletion_by_university"
	CodeNoRequest       = "deletion_not_requested"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrSelfDeletionOff):
		apierror.Write(w, http.StatusForbidden, CodeSelfDeletionOff)
	case errors.Is(err, ErrNoRequest):
		apierror.Write(w, http.StatusNotFound, CodeNoRequest)
	default:
		slog.Error("privacy request failed", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
	}
}
