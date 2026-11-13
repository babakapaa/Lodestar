package sandbox

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const (
	CodeInvalidClass    = "invalid_gpu_class"
	CodeClassNotFound   = "gpu_class_not_found"
	CodeClassTaken      = "gpu_class_taken"
	CodeInvalidTerm     = "invalid_term"
	CodeTermNotFound    = "term_not_found"
	CodeTermOverlap     = "term_overlap"
	CodeNoTerm          = "no_term"
	CodeInvalidTransfer = "invalid_transfer"
	CodeNotInTeam       = "not_in_team"
	CodeNotEnoughHours  = "not_enough_gpu_hours"
	CodeInvalidTopup    = "invalid_topup"
	CodeWorkNotFound    = "sandbox_work_not_found"
	CodeTopupNotFound   = "topup_not_found"
	CodeTopupPending    = "topup_pending"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrInvalidClass, http.StatusBadRequest, CodeInvalidClass},
	{ErrClassNotFound, http.StatusNotFound, CodeClassNotFound},
	{ErrClassTaken, http.StatusConflict, CodeClassTaken},
	{ErrInvalidTerm, http.StatusBadRequest, CodeInvalidTerm},
	{ErrTermNotFound, http.StatusNotFound, CodeTermNotFound},
	{ErrTermOverlap, http.StatusConflict, CodeTermOverlap},
	{ErrNoTerm, http.StatusConflict, CodeNoTerm},
	{ErrInvalidTransfer, http.StatusBadRequest, CodeInvalidTransfer},
	{ErrNotInTeam, http.StatusForbidden, CodeNotInTeam},
	{ErrNotEnoughHours, http.StatusConflict, CodeNotEnoughHours},
	{ErrInvalidTopup, http.StatusBadRequest, CodeInvalidTopup},
	{ErrWorkNotFound, http.StatusNotFound, CodeWorkNotFound},
	{ErrTopupNotFound, http.StatusNotFound, CodeTopupNotFound},
	{ErrTopupPending, http.StatusConflict, CodeTopupPending},
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

func respond(w http.ResponseWriter, status int, body any, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, status, body)
}
