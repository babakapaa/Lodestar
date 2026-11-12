package study

type LessonStatus string

const (
	LessonPlanned   LessonStatus = "planned"
	LessonChanged   LessonStatus = "changed"
	LessonCancelled LessonStatus = "cancelled"
)
