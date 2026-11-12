package access

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeForbidden          = "forbidden"
	CodeRoleHierarchy      = "role_hierarchy"
	CodePermissionNotHeld  = "permission_not_held"
	CodeUnknownPermission  = "unknown_permission"
	CodeRoleNotFound       = "role_not_found"
	CodeAssignmentNotFound = "assignment_not_found"
	CodeUserNotFound       = "user_not_found"
	CodeInvalidRole        = "invalid_role"
	CodeInvalidScope       = "invalid_scope"
	CodeInvalidPosition    = "invalid_position"
	CodeDefaultRole        = "default_role"
	CodeInvalidRequest     = "invalid_request"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrForbidden, http.StatusForbidden, CodeForbidden},
	{ErrHierarchy, http.StatusForbidden, CodeRoleHierarchy},
	{ErrPermissionNotHeld, http.StatusForbidden, CodePermissionNotHeld},
	{ErrUnknownPermission, http.StatusBadRequest, CodeUnknownPermission},
	{ErrRoleNotFound, http.StatusNotFound, CodeRoleNotFound},
	{ErrAssignmentNotFound, http.StatusNotFound, CodeAssignmentNotFound},
	{ErrUserNotFound, http.StatusNotFound, CodeUserNotFound},
	{ErrInvalidRole, http.StatusBadRequest, CodeInvalidRole},
	{ErrInvalidScope, http.StatusBadRequest, CodeInvalidScope},
	{ErrInvalidPosition, http.StatusBadRequest, CodeInvalidPosition},
	{ErrDefaultRole, http.StatusBadRequest, CodeDefaultRole},
}

func writeError(w http.ResponseWriter, err error) {
	for _, e := range errorCodes {
		if errors.Is(err, e.err) {
			apierror.Write(w, e.status, e.code)
			return
		}
	}
	slog.Error("access request failed", "error", err)
	apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
}

func badRequest(w http.ResponseWriter) {
	apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
}
