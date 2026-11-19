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
	Weights       *GradeWeights  `json:"gradeWeights"`
	SearchMap     *MapRules      `json:"searchMap"`
	NavigatorMode *string        `json:"navigatorMode"`
	NavigatorHelp bool           `json:"navigatorSummary"`
	Defense       *DefenseFormat `json:"defense"`
	Team          *TeamRules     `json:"team"`
	PublishedAt   *time.Time     `json:"publishedAt"`
	CreatedAt     time.Time      `json:"createdAt"`
}

func (a Assignment) Published() bool {
	return a.PublishedAt != nil
}

func (a Assignment) maxScore() float64 {
	if a.Grading == GradingPass {
		return 1
	}
	if a.MaxPoints == nil {
		return 0
	}
	return *a.MaxPoints
}
