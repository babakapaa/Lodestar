package team

import (
	"strings"
	"unicode/utf8"
)

const maxLabel = 40

type LabelInput struct {
	Label string `json:"label"`
}

func (in LabelInput) normalized() (LabelInput, error) {
	in.Label = strings.TrimSpace(in.Label)
	if utf8.RuneCountInString(in.Label) > maxLabel {
		return in, ErrInvalidTeam
	}
	return in, nil
}
