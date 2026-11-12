package notes

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeNoteNotFound    = "note_not_found"
	CodeInvalidNote     = "invalid_note"
	CodeNoteExists      = "note_exists"
	CodeReadOnly        = "note_read_only"
	CodeNoteDeleted     = "note_deleted"
	CodeInvalidUpdate   = "invalid_note_update"
	CodeStaleSnapshot   = "stale_snapshot"
	CodeQuotaExceeded   = "notes_quota_exceeded"
	CodeUserNotFound    = "user_not_found"
	CodeNotOwner        = "not_note_owner"
	CodeVersionNotFound = "note_version_not_found"
	CodeInvalidSource   = "invalid_source"
	CodeSourceNotFound  = "source_not_found"
	CodeLookupFailed    = "source_lookup_failed"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrNoteNotFound, http.StatusNotFound, CodeNoteNotFound},
	{ErrInvalidNote, http.StatusBadRequest, CodeInvalidNote},
	{ErrNoteExists, http.StatusConflict, CodeNoteExists},
	{ErrReadOnly, http.StatusForbidden, CodeReadOnly},
	{ErrNoteDeleted, http.StatusConflict, CodeNoteDeleted},
	{ErrInvalidUpdate, http.StatusBadRequest, CodeInvalidUpdate},
	{ErrStaleSnapshot, http.StatusConflict, CodeStaleSnapshot},
	{ErrQuotaExceeded, http.StatusRequestEntityTooLarge, CodeQuotaExceeded},
	{ErrUserNotFound, http.StatusNotFound, CodeUserNotFound},
	{ErrNotOwner, http.StatusForbidden, CodeNotOwner},
	{ErrVersionNotFound, http.StatusNotFound, CodeVersionNotFound},
	{ErrInvalidSource, http.StatusBadRequest, CodeInvalidSource},
	{ErrSourceNotFound, http.StatusNotFound, CodeSourceNotFound},
	{ErrLookupFailed, http.StatusUnprocessableEntity, CodeLookupFailed},
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
