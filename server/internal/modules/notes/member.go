package notes

import "time"

type Member struct {
	UserID      string    `json:"userId"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	CanEdit     bool      `json:"canEdit"`
	AddedAt     time.Time `json:"addedAt"`
}
