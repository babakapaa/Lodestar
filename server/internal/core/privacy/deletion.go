package privacy

import "time"

type Deletion struct {
	RequestedAt time.Time `json:"requestedAt"`
	DueAt       time.Time `json:"dueAt"`
}
