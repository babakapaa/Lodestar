package notes

import "time"

type sourceRow struct {
	ID        string
	OwnerID   string
	Data      []byte
	FileID    *string
	CreatedAt time.Time
	ChangedAt time.Time
}
