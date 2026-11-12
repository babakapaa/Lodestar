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
}

func (in SubmissionInput) normalized() (SubmissionInput, error) {
	in.Body = strings.TrimSpace(in.Body)
	if in.ClientID == "" || len(in.Files) > maxAttachments || (in.Body == "" && len(in.Files) == 0) {
		return in, ErrInvalidSubmission
	}
	return in, nil
}
