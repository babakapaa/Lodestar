package team

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

type Member struct {
	identity.Person
	Label    string    `json:"label"`
	JoinedAt time.Time `json:"joinedAt"`
}
