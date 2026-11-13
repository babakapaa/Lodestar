package lti

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
)

const (
	CodePlatformNotFound   = "lti_platform_not_found"
	CodeInvalidPlatform    = "invalid_lti_platform"
	CodePlatformTaken      = "lti_platform_taken"
	CodePlatformAmbiguous  = "lti_platform_ambiguous"
	CodeContextNotFound    = "lti_context_not_found"
	CodeUnknownLinkKind    = "unknown_link_kind"
	CodeLinkTargetNotFound = "link_target_not_found"
	CodeInvalidState       = "invalid_state"
	CodeTenantNotFound     = "tenant_not_found"
	CodeNotLinked          = "lti_context_not_linked"
	CodePlatformFailed     = "lms_unavailable"
	CodeSourceFixed        = "assignment_source_fixed"
	CodeInvalidSource      = "invalid_assignment_source"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrPlatformNotFound, http.StatusNotFound, CodePlatformNotFound},
	{ErrInvalidPlatform, http.StatusBadRequest, CodeInvalidPlatform},
	{ErrPlatformTaken, http.StatusConflict, CodePlatformTaken},
	{ErrPlatformAmbiguous, http.StatusBadRequest, CodePlatformAmbiguous},
	{ErrContextNotFound, http.StatusNotFound, CodeContextNotFound},
	{ErrUnknownLinkKind, http.StatusBadRequest, CodeUnknownLinkKind},
	{module.ErrLinkTargetNotFound, http.StatusNotFound, CodeLinkTargetNotFound},
	{ErrInvalidState, http.StatusBadRequest, CodeInvalidState},
	{tenant.ErrNotFound, http.StatusNotFound, CodeTenantNotFound},
	{ErrNotLinked, http.StatusConflict, CodeNotLinked},
	{ErrPlatformUnavailable, http.StatusBadGateway, CodePlatformFailed},
	{ErrSourceFixed, http.StatusConflict, CodeSourceFixed},
	{ErrInvalidSource, http.StatusBadRequest, CodeInvalidSource},
}

func writeError(w http.ResponseWriter, err error) {
	for _, e := range errorCodes {
		if errors.Is(err, e.err) {
			apierror.Write(w, e.status, e.code)
			return
		}
	}
	access.WriteError(w, err)
}

func badRequest(w http.ResponseWriter) {
	apierror.Write(w, http.StatusBadRequest, access.CodeInvalidRequest)
}
