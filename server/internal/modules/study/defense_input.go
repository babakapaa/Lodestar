package study

import "strings"

const maxSummary = 20000

type DefenseInput struct {
	Questions []DefenseQuestion `json:"questions"`
	Summary   string            `json:"summary"`
}

func (in DefenseInput) normalized() (DefenseInput, error) {
	in.Summary = strings.TrimSpace(in.Summary)
	if len(in.Questions) > maxQuestions || len(in.Summary) > maxSummary {
		return in, ErrInvalidDefense
	}
	out := make([]DefenseQuestion, 0, len(in.Questions))
	for _, q := range in.Questions {
		q, ok := q.normalized()
		if !ok {
			return in, ErrInvalidDefense
		}
		out = append(out, q)
	}
	in.Questions = out
	return in, nil
}
