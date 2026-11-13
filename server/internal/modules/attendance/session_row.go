package attendance

import "time"

type sessionRow struct {
	ID          string
	SeriesID    string
	Date        string
	CourseID    string
	CourseTitle string
	Room        string
	StartsAt    time.Time
	EndsAt      time.Time
	Secret      []byte
	Code        string
	CodeUntil   time.Time
	ClosedAt    *time.Time
}

func (r sessionRow) open() bool {
	return r.ClosedAt == nil
}
