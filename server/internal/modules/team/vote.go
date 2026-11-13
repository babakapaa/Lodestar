package team

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

const VoteWaiting = "waiting"

type Vote struct {
	identity.Person
	Decision string     `json:"decision"`
	Comment  string     `json:"comment"`
	At       *time.Time `json:"at"`
}
