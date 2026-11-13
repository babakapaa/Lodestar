package chat

import "time"

type Member struct {
	UserID      string    `json:"userId"`
	DisplayName string    `json:"displayName"`
	Role        string    `json:"role"`
	Synced      bool      `json:"synced"`
	JoinedAt    time.Time `json:"joinedAt"`
}
