package sandbox

import "time"

type Template struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	About     string    `json:"about"`
	Image     string    `json:"image"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"createdAt"`
}
