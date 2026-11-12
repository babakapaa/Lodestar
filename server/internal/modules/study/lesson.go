package study

import "time"

type Lesson struct {
	SeriesID    string       `json:"seriesId"`
	CourseID    string       `json:"courseId"`
	CourseTitle string       `json:"courseTitle"`
	UnitID      *string      `json:"unitId"`
	Kind        string       `json:"kind"`
	Date        Day          `json:"date"`
	StartsAt    time.Time    `json:"startsAt"`
	EndsAt      time.Time    `json:"endsAt"`
	Room        string       `json:"room"`
	TeacherID   *string      `json:"teacherId"`
	Status      LessonStatus `json:"status"`
	Note        string       `json:"note"`
}
