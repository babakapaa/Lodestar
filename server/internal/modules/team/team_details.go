package team

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type TeamDetails struct {
	Team
	Scope          Scope             `json:"scope"`
	Manage         bool              `json:"manage"`
	InviteCode     *string           `json:"inviteCode"`
	Invited        []identity.Person `json:"invited"`
	History        []HistoryEntry    `json:"history"`
	Request        *LeaveRequest     `json:"request"`
	CanLeaveFreely bool              `json:"canLeaveFreely"`
}
