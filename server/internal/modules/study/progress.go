package study

import "time"

type Progress struct {
	SubmittedAt *time.Time `json:"submittedAt"`
	Grade       *Grade     `json:"grade"`
}
