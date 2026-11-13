package navigator

import "time"

type Material struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"courseId"`
	Title     string    `json:"title"`
	Reference string    `json:"reference"`
	URL       string    `json:"url"`
	Passages  int       `json:"passages"`
	Indexed   int       `json:"indexed"`
	CreatedAt time.Time `json:"createdAt"`
}
