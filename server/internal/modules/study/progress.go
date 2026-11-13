package study

import "time"

type Progress struct {
	SubmittedAt *time.Time `json:"submittedAt"`
	Grade       *Grade     `json:"grade"`
}

func (p *Progress) rewind(at time.Time) {
	if p.SubmittedAt != nil && p.SubmittedAt.After(at) {
		p.SubmittedAt = nil
	}
	if p.Grade != nil && p.Grade.GradedAt.After(at) {
		p.Grade = nil
	}
}
