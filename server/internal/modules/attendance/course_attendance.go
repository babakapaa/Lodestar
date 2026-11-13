package attendance

type CourseAttendance struct {
	CourseID    string          `json:"courseId"`
	CourseTitle string          `json:"courseTitle"`
	From        string          `json:"from"`
	To          string          `json:"to"`
	AbsentAlert int             `json:"absentAlert"`
	Lessons     []LessonSummary `json:"lessons"`
	Students    []StudentRecord `json:"students"`
}
