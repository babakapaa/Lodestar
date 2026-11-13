package attendance

import "time"

type CheckResult struct {
	SessionID   string     `json:"sessionId"`
	CourseTitle string     `json:"courseTitle"`
	StartsAt    time.Time  `json:"startsAt"`
	Status      MarkStatus `json:"status"`
}
