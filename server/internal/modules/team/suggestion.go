package team

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type Suggestion struct {
	Teams   []SuggestedTeam   `json:"teams"`
	Without []identity.Person `json:"without"`
}
