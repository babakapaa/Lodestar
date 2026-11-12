package notes

import "time"

type noteRow struct {
	ID        string
	OwnerID   string
	Role      Role
	Meta      []byte
	Seq       int64
	BaseSeq   int64
	Bytes     int64
	CreatedAt time.Time
	ChangedAt time.Time
	DeletedAt *time.Time
}
