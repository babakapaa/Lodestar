package league

import (
	"strings"
	"unicode/utf8"
)

type AnswerInput struct {
	Accept bool   `json:"accept"`
	Answer string `json:"answer"`
	Points int    `json:"points"`
}

func (in AnswerInput) normalized() (AnswerInput, error) {
	in.Answer = strings.TrimSpace(in.Answer)
	if in.Answer == "" || utf8.RuneCountInString(in.Answer) > 1000 || in.Points < -maxWeight || in.Points > maxWeight || (!in.Accept && in.Points != 0) {
		return in, ErrInvalidInput
	}
	return in, nil
}
