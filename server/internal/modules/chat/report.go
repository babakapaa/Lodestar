package chat

import "time"

type Report struct {
	ID         string            `json:"id"`
	ChatID     *string           `json:"chatId"`
	ChatTitle  string            `json:"chatTitle"`
	ReporterID *string           `json:"reporterId"`
	Reason     string            `json:"reason"`
	Messages   []ReportedMessage `json:"messages"`
	CreatedAt  time.Time         `json:"createdAt"`
	ResolvedAt *time.Time        `json:"resolvedAt"`
}
