package chat

import (
	"strings"
	"unicode/utf8"
)

const (
	maxReason   = 1000
	maxReported = 20
)

type ReportInput struct {
	ChatID     string   `json:"chatId"`
	MessageIDs []string `json:"messageIds"`
	Reason     string   `json:"reason"`
}

func (in ReportInput) normalized() (ReportInput, error) {
	in.Reason = strings.TrimSpace(in.Reason)
	if !validID(in.ChatID) || in.Reason == "" || utf8.RuneCountInString(in.Reason) > maxReason {
		return in, ErrInvalidReport
	}
	if len(in.MessageIDs) == 0 || len(in.MessageIDs) > maxReported || !allValid(in.MessageIDs) {
		return in, ErrInvalidReport
	}
	return in, nil
}
