package league

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const (
	CodeLeagueNotFound     = "league_not_found"
	CodeInvalidRules       = "invalid_league_rules"
	CodeInvalidInput       = "invalid_league_input"
	CodePlayerNotFound     = "league_player_not_found"
	CodeTransitionNotFound = "league_move_not_found"
	CodeDisputeNotFound    = "league_dispute_not_found"
	CodeDisputeTarget      = "league_dispute_target"
	CodeDisputeOpen        = "league_dispute_open"
	CodeLeagueArchived     = "league_archived"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrLeagueNotFound, http.StatusNotFound, CodeLeagueNotFound},
	{ErrInvalidRules, http.StatusBadRequest, CodeInvalidRules},
	{ErrInvalidInput, http.StatusBadRequest, CodeInvalidInput},
	{ErrPlayerNotFound, http.StatusNotFound, CodePlayerNotFound},
	{ErrTransitionNotFound, http.StatusNotFound, CodeTransitionNotFound},
	{ErrDisputeNotFound, http.StatusNotFound, CodeDisputeNotFound},
	{ErrDisputeTarget, http.StatusBadRequest, CodeDisputeTarget},
	{ErrDisputeOpen, http.StatusConflict, CodeDisputeOpen},
	{ErrLeagueArchived, http.StatusConflict, CodeLeagueArchived},
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
