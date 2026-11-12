package notes

import "time"

type Source struct {
	ID string `json:"id"`
	SourceFields
	FileID    *string   `json:"fileId"`
	CreatedAt time.Time `json:"createdAt"`
	ChangedAt time.Time `json:"changedAt"`
}
