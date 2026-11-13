package chat

import (
	"strings"
	"unicode/utf8"
)

const maxText = 4000

type MessageInput struct {
	ClientID string `json:"clientId"`
	Text     string `json:"text"`
	ReplyTo  *int64 `json:"replyTo"`
}

func (in MessageInput) normalized() (MessageInput, error) {
	in.ClientID = strings.TrimSpace(in.ClientID)
	text, err := normalizedText(in.Text)
	in.Text = text
	if err != nil || in.ClientID == "" || len(in.ClientID) > 100 || (in.ReplyTo != nil && *in.ReplyTo < 1) {
		return in, ErrInvalidMessage
	}
	return in, nil
}

func normalizedText(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" || utf8.RuneCountInString(text) > maxText {
		return text, ErrInvalidMessage
	}
	return text, nil
}
