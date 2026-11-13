package module

type LessonCourse struct {
	ID      string
	Title   string
	Leads   []string
	Members []string
	Lessons []LessonRef
}
