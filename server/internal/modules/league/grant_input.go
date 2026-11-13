package league

import (
	"strings"
	"unicode/utf8"
)

type GrantInput struct {
	AwardID string `json:"awardId"`
	UserID  string `json:"userId"`
	Reason  string `json:"reason"`
}

func (in GrantInput) normalized() (GrantInput, error) {
	in.Reason = strings.TrimSpace(in.Reason)
	if in.AwardID == "" || in.UserID == "" || in.Reason == "" || utf8.RuneCountInString(in.Reason) > 1000 {
		return in, ErrInvalidInput
	}
	return in, nil
}
