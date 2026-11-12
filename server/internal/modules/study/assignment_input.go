package study

import (
	"strings"
	"time"
)

const maxAttachments = 20

type AssignmentInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueAt       *time.Time `json:"dueAt"`
	Grading     Grading    `json:"grading"`
	MaxPoints   *float64   `json:"maxPoints"`
	Files       []string   `json:"files"`
	Published   bool       `json:"published"`
}

func (in AssignmentInput) normalized() (AssignmentInput, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)
	if in.Grading == "" {
		in.Grading = GradingNone
	}
	if !in.Grading.valid() || len(in.Files) > maxAttachments {
		return in, ErrInvalidAssignment
	}
	if (in.Grading == GradingPoints) != (in.MaxPoints != nil) {
		return in, ErrInvalidAssignment
	}
	return in, nil
}
