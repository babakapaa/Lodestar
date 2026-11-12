package structure

import "time"

type Unit struct {
	ID         string     `json:"id"`
	ParentID   *string    `json:"parentId"`
	Kind       string     `json:"kind"`
	Name       string     `json:"name"`
	Code       *string    `json:"code"`
	ArchivedAt *time.Time `json:"archivedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}
