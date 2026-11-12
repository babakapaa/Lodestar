package study

import "time"

type Grade struct {
	AssignmentID string    `json:"assignmentId"`
	StudentID    string    `json:"studentId"`
	Points       *float64  `json:"points"`
	Passed       *bool     `json:"passed"`
	Comment      string    `json:"comment"`
	GradedBy     *string   `json:"gradedBy"`
	GradedAt     time.Time `json:"gradedAt"`
}
