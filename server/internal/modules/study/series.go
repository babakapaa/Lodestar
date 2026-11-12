package study

type Series struct {
	ID         string  `json:"id"`
	CourseID   string  `json:"courseId"`
	UnitID     *string `json:"unitId"`
	TeacherID  *string `json:"teacherId"`
	Kind       string  `json:"kind"`
	Room       string  `json:"room"`
	StartsOn   Day     `json:"startsOn"`
	EndsOn     Day     `json:"endsOn"`
	Weekday    int     `json:"weekday"`
	StartTime  Clock   `json:"startTime"`
	Minutes    int     `json:"minutes"`
	EveryWeeks int     `json:"everyWeeks"`
}
