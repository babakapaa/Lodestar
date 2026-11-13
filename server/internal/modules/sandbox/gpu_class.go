package sandbox

import "time"

type GPUClass struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	MemoryGB  int       `json:"memoryGb"`
	Rate      float64   `json:"rate"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"createdAt"`
}
