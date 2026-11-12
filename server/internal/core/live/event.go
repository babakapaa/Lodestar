package live

import (
	"encoding/json"
	"time"
)

const ResyncType = "core.resync"

type Event struct {
	ID      int64           `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	At      time.Time       `json:"at"`
}
