package structure

import (
	"encoding/json"
	"time"
)

const (
	ActionCreated       = "created"
	ActionUpdated       = "updated"
	ActionMoved         = "moved"
	ActionArchived      = "archived"
	ActionRestored      = "restored"
	ActionMemberAdded   = "member_added"
	ActionMemberRemoved = "member_removed"
)

type Event struct {
	ID      int64           `json:"id"`
	ActorID *string         `json:"actorId"`
	Action  string          `json:"action"`
	Data    json.RawMessage `json:"data"`
	At      time.Time       `json:"at"`
}
