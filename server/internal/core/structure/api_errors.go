package structure

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeUnitNotFound = "unit_not_found"
	CodeInvalidUnit  = "invalid_unit"
	CodeCodeTaken    = "unit_code_taken"
	CodeCycle        = "unit_cycle"
	CodeArchived     = "unit_archived"
	CodeUserNotFound = "user_not_found"
	CodeBadImport    = "bad_import"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrUnitNotFound, http.StatusNotFound, CodeUnitNotFound},
	{ErrInvalidUnit, http.StatusBadRequest, CodeInvalidUnit},
	{ErrCodeTaken, http.StatusConflict, CodeCodeTaken},
	{ErrCycle, http.StatusBadRequest, CodeCycle},
	{ErrArchived, http.StatusConflict, CodeArchived},
	{ErrUserNotFound, http.StatusNotFound, CodeUserNotFound},
	{ErrBadImport, http.StatusBadRequest, CodeBadImport},
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
