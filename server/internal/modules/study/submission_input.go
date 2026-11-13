package study

import (
	"strings"
	"time"
)

type SubmissionInput struct {
	ClientID   string     `json:"clientId"`
	Body       string     `json:"body"`
	Files      []string   `json:"files"`
	ComposedAt *time.Time `json:"composedAt"`
	SearchMap  *SearchMap `json:"searchMap"`
}

func (in SubmissionInput) normalized() (SubmissionInput, error) {
	in.Body = strings.TrimSpace(in.Body)
	if in.ClientID == "" || len(in.Files) > maxAttachments || (in.Body == "" && len(in.Files) == 0) {
		return in, ErrInvalidSubmission
	}
	if in.SearchMap != nil {
		m, err := in.SearchMap.normalized()
		if err != nil {
			return in, err
		}
		in.SearchMap = &m
	}
	return in, nil
}
