package league

import (
	"strings"
	"unicode/utf8"
)

type DisputeInput struct {
	EntryID      *string `json:"entryId"`
	TransitionID *string `json:"transitionId"`
	Text         string  `json:"text"`
}

func (in DisputeInput) normalized() (DisputeInput, error) {
	in.Text = strings.TrimSpace(in.Text)
	if (in.EntryID == nil) == (in.TransitionID == nil) || in.Text == "" || utf8.RuneCountInString(in.Text) > 2000 {
		return in, ErrInvalidInput
	}
	return in, nil
}
