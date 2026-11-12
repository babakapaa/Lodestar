package search

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeInvalidQuery = "invalid_query"
	CodeUnknownKind  = "unknown_search_kind"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidQuery):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidQuery)
	case errors.Is(err, ErrUnknownKind):
		apierror.Write(w, http.StatusBadRequest, CodeUnknownKind)
	default:
		slog.Error("search failed", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
	}
}
