package notes

import "time"

type Update struct {
	Seq       int64     `json:"seq"`
	Key       string    `json:"key"`
	AuthorID  *string   `json:"authorId"`
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"createdAt"`
}
