package attendance

import "time"

type Mark struct {
	UserID   string      `json:"userId"`
	Name     string      `json:"name"`
	Status   *MarkStatus `json:"status"`
	Method   *MarkMethod `json:"method"`
	MarkedAt *time.Time  `json:"markedAt"`
}
