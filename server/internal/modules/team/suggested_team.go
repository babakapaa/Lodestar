package team

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type SuggestedTeam struct {
	Members []identity.Person `json:"members"`
	Why     string            `json:"why"`
}
