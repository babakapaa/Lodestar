package audit

import (
	"encoding/json"
	"time"
)

type Entry struct {
	ID         int64           `json:"id"`
	ActorID    *string         `json:"actorId"`
	Action     string          `json:"action"`
	TargetKind string          `json:"targetKind"`
	TargetID   string          `json:"targetId"`
	Details    json.RawMessage `json:"details"`
	At         time.Time       `json:"at"`
}
