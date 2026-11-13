package chat

import "time"

type ReportedMessage struct {
	Seq        int64     `json:"seq"`
	AuthorID   *string   `json:"authorId"`
	AuthorName string    `json:"authorName"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"createdAt"`
}
