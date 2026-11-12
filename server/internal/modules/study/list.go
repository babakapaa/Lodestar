package study

type List string

const (
	ListTeachers List = "teachers"
	ListUnits    List = "units"
	ListStudents List = "students"
)

func (l List) table() string {
	return "course_" + string(l)
}

func (l List) column() string {
	if l == ListUnits {
		return "unit_id"
	}
	return "user_id"
}

func (l List) missing() error {
	if l == ListUnits {
		return ErrUnitNotFound
	}
	return ErrUserNotFound
}

func (l List) actions() (string, string) {
	switch l {
	case ListTeachers:
		return ActionTeacherAdded, ActionTeacherRemoved
	case ListUnits:
		return ActionUnitAdded, ActionUnitRemoved
	default:
		return ActionStudentAdded, ActionStudentRemoved
	}
}
