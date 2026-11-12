package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
)

const (
	CodeUnauthorized     = "unauthorized"
	CodeInvalidRequest   = "invalid_request"
	CodeInvalidGrant     = "invalid_grant"
	CodeTenantNotFound   = "tenant_not_found"
	CodeProviderNotFound = "provider_not_found"
	CodeRedirectRejected = "redirect_not_allowed"
	CodeInvalidState     = "invalid_state"
	CodeSessionNotFound  = "session_not_found"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, tenant.ErrNotFound):
		apierror.Write(w, http.StatusNotFound, CodeTenantNotFound)
	case errors.Is(err, ErrProviderNotFound):
		apierror.Write(w, http.StatusNotFound, CodeProviderNotFound)
	case errors.Is(err, ErrSessionNotFound):
		apierror.Write(w, http.StatusNotFound, CodeSessionNotFound)
	case errors.Is(err, ErrRedirectNotAllowed):
		apierror.Write(w, http.StatusBadRequest, CodeRedirectRejected)
	case errors.Is(err, ErrInvalidChallenge):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
	case errors.Is(err, ErrInvalidState):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidState)
	case errors.Is(err, identity.ErrNotFound):
		apierror.Write(w, http.StatusUnauthorized, CodeUnauthorized)
	case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrTokenReused):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidGrant)
	default:
		slog.Error("auth request failed", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
	}
}
