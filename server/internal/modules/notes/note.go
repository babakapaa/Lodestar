package notes

import "time"

type Note struct {
	ID        string     `json:"id"`
	OwnerID   string     `json:"ownerId"`
	Role      Role       `json:"role"`
	Meta      Meta       `json:"meta"`
	Seq       int64      `json:"seq"`
	Bytes     int64      `json:"bytes"`
	CreatedAt time.Time  `json:"createdAt"`
	ChangedAt time.Time  `json:"changedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
