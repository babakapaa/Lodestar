package league

import (
	"strings"
	"unicode/utf8"
)

type ReviewInput struct {
	Points  int      `json:"points"`
	Passed  []string `json:"passed"`
	Comment string   `json:"comment"`
}

func (in ReviewInput) normalized(s Solution, task Task) (ReviewInput, error) {
	in.Comment = strings.TrimSpace(in.Comment)
	if in.Points < 0 || in.Points > task.Points || utf8.RuneCountInString(in.Comment) > 2000 {
		return in, ErrInvalidInput
	}
	seen := map[string]bool{}
	for _, id := range in.Passed {
		if !s.has(id) || seen[id] {
			return in, ErrInvalidInput
		}
		seen[id] = true
	}
	if in.Passed == nil {
		in.Passed = []string{}
	}
	return in, nil
}
