package settings

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	CodeUnknownKey     = "unknown_key"
	CodeInvalidValue   = "invalid_value"
	CodeUnknownModule  = "module_not_found"
	CodeRequiredModule = "module_required"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, schema.ErrUnknownKey):
		apierror.Write(w, http.StatusNotFound, CodeUnknownKey)
	case errors.Is(err, schema.ErrInvalidValue):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidValue)
	case errors.Is(err, ErrUnknownModule):
		apierror.Write(w, http.StatusNotFound, CodeUnknownModule)
	case errors.Is(err, ErrRequiredModule):
		apierror.Write(w, http.StatusBadRequest, CodeRequiredModule)
	default:
		access.WriteError(w, err)
	}
}

func badRequest(w http.ResponseWriter) {
	apierror.Write(w, http.StatusBadRequest, access.CodeInvalidRequest)
}
