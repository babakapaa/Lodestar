package study

import "strings"

const (
	QuestionUnderstanding = "understanding"
	QuestionAnalysis      = "analysis"
	QuestionApplication   = "application"
	QuestionPath          = "path"
)

const (
	maxQuestions    = 20
	maxQuestionText = 2000
	maxMark         = 3
)

type DefenseQuestion struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
	Mark *int   `json:"mark"`
	Note string `json:"note"`
}

func (q DefenseQuestion) normalized() (DefenseQuestion, bool) {
	q.Text = strings.TrimSpace(q.Text)
	q.Note = strings.TrimSpace(q.Note)
	switch q.Kind {
	case QuestionUnderstanding, QuestionAnalysis, QuestionApplication, QuestionPath:
	default:
		return q, false
	}
	if q.Text == "" || len(q.Text) > maxQuestionText || len(q.Note) > maxQuestionText {
		return q, false
	}
	if q.Mark != nil && (*q.Mark < 0 || *q.Mark > maxMark) {
		return q, false
	}
	return q, true
}
