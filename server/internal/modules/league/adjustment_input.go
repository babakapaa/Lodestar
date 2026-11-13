package league

import (
	"strings"
	"unicode/utf8"
)

type AdjustmentInput struct {
	UserID string `json:"userId"`
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

func (in AdjustmentInput) normalized() (AdjustmentInput, error) {
	in.Reason = strings.TrimSpace(in.Reason)
	if in.Points == 0 || in.Points < -maxWeight || in.Points > maxWeight || in.Reason == "" || utf8.RuneCountInString(in.Reason) > 1000 {
		return in, ErrInvalidInput
	}
	return in, nil
}
