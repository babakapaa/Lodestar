package navigator

import (
	"strings"
	"unicode/utf8"
)

const maxQuestion = 8000

type MessageInput struct {
	Key    string `json:"key"`
	Text   string `json:"text"`
	Intent Intent `json:"intent"`
	More   bool   `json:"more"`
}

func (in MessageInput) normalized() (MessageInput, error) {
	in.Text = strings.TrimSpace(in.Text)
	in.Key = strings.TrimSpace(in.Key)
	if in.Key == "" || len(in.Key) > 100 || utf8.RuneCountInString(in.Text) > maxQuestion {
		return in, ErrInvalidMessage
	}
	if in.Intent != "" && !in.Intent.Choosable() {
		return in, ErrInvalidMessage
	}
	if in.Text == "" && !in.More {
		return in, ErrInvalidMessage
	}
	return in, nil
}
