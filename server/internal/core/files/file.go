package files

import "time"

const (
	StatusPending = "pending"
	StatusReady   = "ready"
)

type File struct {
	ID          string     `json:"id"`
	OwnerID     string     `json:"ownerId"`
	Name        string     `json:"name"`
	ContentType string     `json:"contentType"`
	Size        int64      `json:"size"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UploadedAt  *time.Time `json:"uploadedAt"`
}
