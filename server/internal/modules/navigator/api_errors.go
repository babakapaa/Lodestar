package navigator

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const (
	CodeUnavailable      = "navigator_unavailable"
	CodeDialogNotFound   = "dialog_not_found"
	CodeInvalidDialog    = "invalid_dialog"
	CodeInvalidMessage   = "invalid_message"
	CodeModelFailed      = "navigator_failed"
	CodeScopeNotFound    = "scope_not_found"
	CodeInvalidMaterial  = "invalid_material"
	CodeMaterialNotFound = "material_not_found"
	CodeJournalNotFound  = "journal_entry_not_found"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrUnavailable, http.StatusServiceUnavailable, CodeUnavailable},
	{ErrDialogNotFound, http.StatusNotFound, CodeDialogNotFound},
	{ErrInvalidDialog, http.StatusBadRequest, CodeInvalidDialog},
	{ErrInvalidMessage, http.StatusBadRequest, CodeInvalidMessage},
	{ErrModelFailed, http.StatusBadGateway, CodeModelFailed},
	{module.ErrBriefNotFound, http.StatusNotFound, CodeScopeNotFound},
	{ErrInvalidMaterial, http.StatusBadRequest, CodeInvalidMaterial},
	{ErrMaterialNotFound, http.StatusNotFound, CodeMaterialNotFound},
	{ErrJournalNotFound, http.StatusNotFound, CodeJournalNotFound},
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
