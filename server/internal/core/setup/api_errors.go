package setup

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
)

const (
	CodeUnknownSection  = "setup_section_not_found"
	CodeUnknownQuestion = "setup_question_not_found"
	CodeWrongSection    = "setup_wrong_section"
	CodeInvalidAnswer   = "invalid_setup_answer"
	CodeInvalidAssignee = "invalid_setup_assignee"
	CodeInvalidStatus   = "setup_section_in_review"
	CodeLegalOpen       = "setup_legal_open"
	CodeNotReady        = "setup_not_ready"
	CodeNoSetup         = "setup_not_found"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrUnknownSection, http.StatusNotFound, CodeUnknownSection},
	{ErrUnknownQuestion, http.StatusNotFound, CodeUnknownQuestion},
	{ErrWrongSection, http.StatusBadRequest, CodeWrongSection},
	{ErrInvalidAnswer, http.StatusBadRequest, CodeInvalidAnswer},
	{ErrInvalidAssignee, http.StatusBadRequest, CodeInvalidAssignee},
	{ErrInvalidStatus, http.StatusConflict, CodeInvalidStatus},
	{ErrLegalOpen, http.StatusConflict, CodeLegalOpen},
	{ErrNotReady, http.StatusConflict, CodeNotReady},
	{ErrNoSetup, http.StatusNotFound, CodeNoSetup},
	{schema.ErrUnknownKey, http.StatusNotFound, settings.CodeUnknownKey},
	{schema.ErrInvalidValue, http.StatusBadRequest, settings.CodeInvalidValue},
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
