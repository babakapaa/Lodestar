package league

import "errors"

var ErrLeagueNotFound = errors.New("league not found or not yours")

var ErrInvalidRules = errors.New("divisions or weights are invalid")

var ErrInvalidInput = errors.New("league input is invalid")

var ErrPlayerNotFound = errors.New("student is not in this league")

var ErrTransitionNotFound = errors.New("transition not found or already decided")

var ErrDisputeNotFound = errors.New("dispute not found or already answered")

var ErrDisputeTarget = errors.New("only own points and applied moves can be disputed")

var ErrDisputeOpen = errors.New("dispute is already open")

var ErrLeagueArchived = errors.New("league is archived with its course")
