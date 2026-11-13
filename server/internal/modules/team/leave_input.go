package team

import (
	"strings"
	"unicode/utf8"
)

const maxReason = 1000

type LeaveInput struct {
	Reason string `json:"reason"`
}

func (in LeaveInput) normalized() (LeaveInput, error) {
	in.Reason = strings.TrimSpace(in.Reason)
	if utf8.RuneCountInString(in.Reason) > maxReason {
		return in, ErrInvalidTeam
	}
	return in, nil
}
