package chat

import "time"

type Call struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chatId"`
	StartedBy *string   `json:"startedBy"`
	StartedAt time.Time `json:"startedAt"`
	People    []string  `json:"people"`
}
