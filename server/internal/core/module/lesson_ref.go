package module

import "time"

type LessonRef struct {
	SeriesID    string
	CourseID    string
	CourseTitle string
	Date        string
	StartsAt    time.Time
	EndsAt      time.Time
	Room        string
	Cancelled   bool
	Leads       []string
	Members     []string
}
