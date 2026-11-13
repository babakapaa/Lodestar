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

var ErrInvalidTask = errors.New("task title, points or team size are invalid")

var ErrTaskNotFound = errors.New("task not found")

var ErrTaskClosed = errors.New("task is closed")

var ErrInvalidSolution = errors.New("solution needs text or links")

var ErrSolutionNotFound = errors.New("solution not found")

var ErrInvalidAward = errors.New("award title, hours or quantity are invalid")

var ErrAwardNotFound = errors.New("award not found")

var ErrAwardGone = errors.New("award is archived or all given out")

var ErrGrantNotFound = errors.New("grant not found")
