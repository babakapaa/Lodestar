package league

import (
	"strings"
	"unicode/utf8"
)

type ExplanationInput struct {
	Text string `json:"text"`
}

func (in ExplanationInput) normalized() (ExplanationInput, error) {
	in.Text = strings.TrimSpace(in.Text)
	if in.Text == "" || utf8.RuneCountInString(in.Text) > 2000 {
		return in, ErrInvalidInput
	}
	return in, nil
}
