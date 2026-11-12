package notify

import (
	"encoding/json"
	"time"
)

const LiveType = "core.notification"

type Notification struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Params    json.RawMessage `json:"params"`
	Link      string          `json:"link"`
	CreatedAt time.Time       `json:"createdAt"`
	ReadAt    *time.Time      `json:"readAt"`
}
