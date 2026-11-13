package attendance

import "time"

type MyMark struct {
	SessionID   string     `json:"sessionId"`
	CourseID    string     `json:"courseId"`
	CourseTitle string     `json:"courseTitle"`
	StartsAt    time.Time  `json:"startsAt"`
	Status      MarkStatus `json:"status"`
	Method      MarkMethod `json:"method"`
}
