package module

type LessonProviding interface {
	LessonSource(k Kernel) LessonSource
}
