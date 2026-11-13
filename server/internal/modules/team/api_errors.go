package team

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const (
	CodeScopeNotFound   = "team_scope_not_found"
	CodeTeamNotFound    = "team_not_found"
	CodeInvalidTeam     = "invalid_team"
	CodeInTeam          = "already_in_team"
	CodeTeamFull        = "team_full"
	CodeNotSelfFormed   = "teams_by_teacher"
	CodeSeparated       = "works_alone"
	CodeScopeArchived   = "teams_archived"
	CodeNotInvited      = "not_invited"
	CodeRequestNotFound = "leave_request_not_found"
	CodeRequestPending  = "leave_request_pending"
	CodeProposalGone    = "proposal_not_found"
	CodeAIUnavailable   = "team_ai_unavailable"
	CodeAIFailed        = "team_ai_failed"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrScopeNotFound, http.StatusNotFound, CodeScopeNotFound},
	{ErrTeamNotFound, http.StatusNotFound, CodeTeamNotFound},
	{ErrInvalidTeam, http.StatusBadRequest, CodeInvalidTeam},
	{ErrInTeam, http.StatusConflict, CodeInTeam},
	{ErrTeamFull, http.StatusConflict, CodeTeamFull},
	{ErrNotSelfFormed, http.StatusForbidden, CodeNotSelfFormed},
	{ErrSeparated, http.StatusConflict, CodeSeparated},
	{ErrScopeArchived, http.StatusConflict, CodeScopeArchived},
	{ErrNotInvited, http.StatusNotFound, CodeNotInvited},
	{ErrRequestNotFound, http.StatusNotFound, CodeRequestNotFound},
	{ErrRequestPending, http.StatusConflict, CodeRequestPending},
	{ErrProposalNotFound, http.StatusNotFound, CodeProposalGone},
	{ErrAIUnavailable, http.StatusServiceUnavailable, CodeAIUnavailable},
	{ErrAIFailed, http.StatusBadGateway, CodeAIFailed},
}

func writeError(w http.ResponseWriter, err error) {
	var problem module.TeamWorkProblem
	if errors.As(err, &problem) {
		apierror.Write(w, problem.Status, problem.Code)
		return
	}
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
