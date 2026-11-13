package study

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/files"
)

type Submission struct {
	ID           string       `json:"id"`
	AssignmentID string       `json:"assignmentId"`
	StudentID    string       `json:"studentId"`
	ClientID     string       `json:"clientId"`
	Body         string       `json:"body"`
	ComposedAt   *time.Time   `json:"composedAt"`
	SearchMap    *SearchMap   `json:"searchMap"`
	SubmittedAt  time.Time    `json:"submittedAt"`
	Files        []files.File `json:"files"`
	Team         *TeamReceipt `json:"team"`
}
