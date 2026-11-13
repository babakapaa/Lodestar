package team

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type TeamCard struct {
	Team      Team             `json:"team"`
	Scope     Scope            `json:"scope"`
	InvitedBy *identity.Person `json:"invitedBy,omitempty"`
}
