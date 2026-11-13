package team

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type Seeker struct {
	identity.Person
	Note string `json:"note"`
}
