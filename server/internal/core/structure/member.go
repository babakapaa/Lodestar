package structure

import "time"

type Member struct {
	UserID      string    `json:"userId"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	JoinedAt    time.Time `json:"joinedAt"`
}
