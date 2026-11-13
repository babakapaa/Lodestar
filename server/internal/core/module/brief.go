package module

import "time"

const (
	ModeStudy   = "study"
	ModeProject = "project"
)

type Brief struct {
	Kind        string     `json:"kind"`
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Course      string     `json:"course"`
	CourseID    string     `json:"courseId"`
	Teaches     bool       `json:"teaches"`
	DueAt       *time.Time `json:"dueAt"`
	Mode        string     `json:"mode"`
}
