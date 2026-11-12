package study

type Agenda struct {
	Date      Day        `json:"date"`
	TimeZone  string     `json:"timeZone"`
	Lessons   []Lesson   `json:"lessons"`
	Deadlines []Deadline `json:"deadlines"`
}
