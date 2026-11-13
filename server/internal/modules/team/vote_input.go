package team

import (
	"strings"
	"unicode/utf8"
)

type VoteInput struct {
	Agree   bool   `json:"agree"`
	Comment string `json:"comment"`
}

func (in VoteInput) normalized() (VoteInput, error) {
	in.Comment = strings.TrimSpace(in.Comment)
	if utf8.RuneCountInString(in.Comment) > maxReason || (!in.Agree && in.Comment == "") {
		return in, ErrInvalidTeam
	}
	return in, nil
}
