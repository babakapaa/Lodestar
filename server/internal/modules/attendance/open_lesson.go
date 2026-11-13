package attendance

import "time"

type OpenLesson struct {
	SessionID   string      `json:"sessionId"`
	CourseTitle string      `json:"courseTitle"`
	Room        string      `json:"room"`
	StartsAt    time.Time   `json:"startsAt"`
	EndsAt      time.Time   `json:"endsAt"`
	Status      *MarkStatus `json:"status"`
	Methods     Methods     `json:"methods"`
}
