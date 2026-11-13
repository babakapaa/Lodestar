package team

import "errors"

var ErrScopeNotFound = errors.New("team scope not found or not yours")

var ErrTeamNotFound = errors.New("team not found")

var ErrInvalidTeam = errors.New("team name, members or label are invalid")

var ErrInTeam = errors.New("student is already in a team here")

var ErrTeamFull = errors.New("team is full")

var ErrNotSelfFormed = errors.New("teams here are formed by the teacher")

var ErrSeparated = errors.New("student works alone by teacher decision")

var ErrScopeArchived = errors.New("teams here are archived")

var ErrNotInvited = errors.New("no invitation to this team")

var ErrRequestNotFound = errors.New("leave request not found")

var ErrRequestPending = errors.New("leave request is already waiting")
