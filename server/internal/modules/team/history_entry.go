package team

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

type HistoryEntry struct {
	Person *identity.Person `json:"person"`
	Event  HistoryEvent     `json:"event"`
	Actor  *identity.Person `json:"actor"`
	At     time.Time        `json:"at"`
}
