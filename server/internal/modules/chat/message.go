package chat

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/files"
)

type Message struct {
	ID        string       `json:"id"`
	ChatID    string       `json:"chatId"`
	Seq       int64        `json:"seq"`
	Rev       int64        `json:"rev"`
	AuthorID  *string      `json:"authorId"`
	ClientID  *string      `json:"clientId"`
	Text      string       `json:"text"`
	ReplyTo   *int64       `json:"replyTo"`
	Mentions  []string     `json:"mentions"`
	Files     []files.File `json:"files"`
	Views     *int         `json:"views"`
	Deleted   bool         `json:"deleted"`
	EditedAt  *time.Time   `json:"editedAt"`
	CreatedAt time.Time    `json:"createdAt"`
}
