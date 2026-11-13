package chat

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeChatNotFound    = "chat_not_found"
	CodeInvalidChat     = "invalid_chat"
	CodeReadOnly        = "chat_read_only"
	CodeInvalidMessage  = "invalid_message"
	CodeMessageNotFound = "message_not_found"
	CodeOriginNotFound  = "origin_not_found"
	CodeSyncedMember    = "synced_member"
	CodeOwnerLeaves     = "owner_leaves"
	CodeBadFile         = "invalid_file"
	CodeBlocked         = "blocked"
	CodeInvalidReport   = "invalid_report"
	CodeReportNotFound  = "report_not_found"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrChatNotFound, http.StatusNotFound, CodeChatNotFound},
	{ErrInvalidChat, http.StatusBadRequest, CodeInvalidChat},
	{ErrReadOnly, http.StatusForbidden, CodeReadOnly},
	{ErrInvalidMessage, http.StatusBadRequest, CodeInvalidMessage},
	{ErrMessageNotFound, http.StatusNotFound, CodeMessageNotFound},
	{ErrOriginNotFound, http.StatusNotFound, CodeOriginNotFound},
	{ErrSyncedMember, http.StatusConflict, CodeSyncedMember},
	{ErrOwnerLeaves, http.StatusConflict, CodeOwnerLeaves},
	{ErrBadFile, http.StatusBadRequest, CodeBadFile},
	{ErrBlocked, http.StatusForbidden, CodeBlocked},
	{ErrInvalidReport, http.StatusBadRequest, CodeInvalidReport},
	{ErrReportNotFound, http.StatusNotFound, CodeReportNotFound},
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
