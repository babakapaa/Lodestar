package notes

import "time"

type Version struct {
	ID        string    `json:"id"`
	Seq       int64     `json:"seq"`
	AuthorID  *string   `json:"authorId"`
	Bytes     int64     `json:"bytes"`
	CreatedAt time.Time `json:"createdAt"`
}
