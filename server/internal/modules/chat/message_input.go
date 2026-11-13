package chat

import (
	"strings"
	"unicode/utf8"
)

const maxText = 4000

const (
	maxFiles    = 10
	maxMentions = 50
)

type MessageInput struct {
	ClientID string   `json:"clientId"`
	Text     string   `json:"text"`
	ReplyTo  *int64   `json:"replyTo"`
	Mentions []string `json:"mentions"`
	FileIDs  []string `json:"fileIds"`
}

func (in MessageInput) normalized() (MessageInput, error) {
	in.ClientID = strings.TrimSpace(in.ClientID)
	text, err := normalizedText(in.Text)
	in.Text = text
	if err != nil && (len(in.FileIDs) == 0 || in.Text != "") {
		return in, ErrInvalidMessage
	}
	if in.ClientID == "" || len(in.ClientID) > 100 || (in.ReplyTo != nil && *in.ReplyTo < 1) {
		return in, ErrInvalidMessage
	}
	if len(in.FileIDs) > maxFiles || len(in.Mentions) > maxMentions || !allValid(in.FileIDs) || !allValid(in.Mentions) {
		return in, ErrInvalidMessage
	}
	return in, nil
}

func allValid(ids []string) bool {
	for _, id := range ids {
		if !validID(id) {
			return false
		}
	}
	return true
}

func normalizedText(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" || utf8.RuneCountInString(text) > maxText {
		return text, ErrInvalidMessage
	}
	return text, nil
}
