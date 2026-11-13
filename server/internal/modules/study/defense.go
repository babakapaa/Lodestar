package study

import "time"

type Defense struct {
	AssignmentID string            `json:"assignmentId"`
	StudentID    string            `json:"studentId"`
	SlotID       *string           `json:"slotId"`
	Selected     bool              `json:"selected"`
	Status       DefenseStatus     `json:"status"`
	Questions    []DefenseQuestion `json:"questions"`
	Summary      string            `json:"summary"`
	ExaminerID   *string           `json:"examinerId"`
	HeldAt       *time.Time        `json:"heldAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}
