package study

import "time"

type Assignment struct {
	ID            string         `json:"id"`
	CourseID      string         `json:"courseId"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	DueAt         *time.Time     `json:"dueAt"`
	Grading       Grading        `json:"grading"`
	MaxPoints     *float64       `json:"maxPoints"`
	SearchMap     *MapRules      `json:"searchMap"`
	NavigatorMode *string        `json:"navigatorMode"`
	Defense       *DefenseFormat `json:"defense"`
	PublishedAt   *time.Time     `json:"publishedAt"`
	CreatedAt     time.Time      `json:"createdAt"`
}

func (a Assignment) Published() bool {
	return a.PublishedAt != nil
}
