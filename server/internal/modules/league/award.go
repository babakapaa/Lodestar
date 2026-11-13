package league

import "time"

type Award struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	About      string     `json:"about"`
	GPUHours   int        `json:"gpuHours"`
	Quantity   *int       `json:"quantity"`
	Granted    int        `json:"granted"`
	ArchivedAt *time.Time `json:"archivedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func (a Award) left() bool {
	return a.ArchivedAt == nil && (a.Quantity == nil || a.Granted < *a.Quantity)
}
