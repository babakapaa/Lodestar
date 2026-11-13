package lti

import "time"

type LineItem struct {
	ID           string     `json:"id,omitempty"`
	Label        string     `json:"label"`
	ScoreMaximum float64    `json:"scoreMaximum"`
	ResourceID   string     `json:"resourceId,omitempty"`
	Tag          string     `json:"tag,omitempty"`
	EndDateTime  *time.Time `json:"endDateTime,omitempty"`
}
