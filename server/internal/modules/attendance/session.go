package attendance

import "time"

type Session struct {
	ID          string     `json:"id"`
	SeriesID    string     `json:"seriesId"`
	Date        string     `json:"date"`
	CourseID    string     `json:"courseId"`
	CourseTitle string     `json:"courseTitle"`
	Room        string     `json:"room"`
	StartsAt    time.Time  `json:"startsAt"`
	EndsAt      time.Time  `json:"endsAt"`
	Open        bool       `json:"open"`
	Code        string     `json:"code,omitempty"`
	CodeUntil   *time.Time `json:"codeUntil,omitempty"`
	Methods     Methods    `json:"methods"`
	Counts      Counts     `json:"counts"`
	Marks       []Mark     `json:"marks"`
}
