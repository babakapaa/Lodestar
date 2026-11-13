package sandbox

import (
	"strings"
	"unicode/utf8"
)

type TopupInput struct {
	WorkKind string  `json:"workKind"`
	WorkID   string  `json:"workId"`
	Hours    float64 `json:"hours"`
	Reason   string  `json:"reason"`
}

func (in TopupInput) normalized() (TopupInput, error) {
	in.Reason = strings.TrimSpace(in.Reason)
	in.Hours = hours(in.Hours)
	if in.WorkKind == "" || in.WorkID == "" || in.Hours <= 0 || in.Hours > 10000 || utf8.RuneCountInString(in.Reason) > 1000 {
		return in, ErrInvalidTopup
	}
	return in, nil
}
