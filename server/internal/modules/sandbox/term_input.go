package sandbox

import (
	"strings"
	"time"
	"unicode/utf8"
)

type TermInput struct {
	Name     string `json:"name"`
	StartsOn string `json:"startsOn"`
	EndsOn   string `json:"endsOn"`
}

func (in TermInput) normalized() (TermInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	from, errFrom := time.Parse(time.DateOnly, in.StartsOn)
	to, errTo := time.Parse(time.DateOnly, in.EndsOn)
	if in.Name == "" || utf8.RuneCountInString(in.Name) > 100 || errFrom != nil || errTo != nil || !to.After(from) {
		return in, ErrInvalidTerm
	}
	return in, nil
}
