package auth

import "time"

type Session struct {
	ID         string    `json:"id"`
	DeviceName string    `json:"deviceName"`
	CreatedAt  time.Time `json:"createdAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
	Current    bool      `json:"current"`
}
