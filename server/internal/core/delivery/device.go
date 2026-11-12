package delivery

import "time"

type Device struct {
	ID         string    `json:"id"`
	Platform   string    `json:"platform"`
	Name       string    `json:"name"`
	Language   string    `json:"language"`
	CreatedAt  time.Time `json:"createdAt"`
	LastSeenAt time.Time `json:"lastSeenAt"`
}
