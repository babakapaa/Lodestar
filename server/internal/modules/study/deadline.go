package study

import "time"

type Deadline struct {
	AssignmentID string     `json:"assignmentId"`
	CourseID     string     `json:"courseId"`
	CourseTitle  string     `json:"courseTitle"`
	Title        string     `json:"title"`
	DueAt        time.Time  `json:"dueAt"`
	SubmittedAt  *time.Time `json:"submittedAt"`
}
