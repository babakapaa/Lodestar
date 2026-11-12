package study

import "strings"

type GradeInput struct {
	Points  *float64 `json:"points"`
	Passed  *bool    `json:"passed"`
	Comment string   `json:"comment"`
}

func (in GradeInput) normalized(a Assignment) (GradeInput, error) {
	in.Comment = strings.TrimSpace(in.Comment)
	switch a.Grading {
	case GradingPoints:
		if in.Points == nil || in.Passed != nil || *in.Points < 0 || *in.Points > *a.MaxPoints {
			return in, ErrInvalidGrade
		}
	case GradingPass:
		if in.Passed == nil || in.Points != nil {
			return in, ErrInvalidGrade
		}
	default:
		if in.Points != nil || in.Passed != nil || in.Comment == "" {
			return in, ErrInvalidGrade
		}
	}
	return in, nil
}
