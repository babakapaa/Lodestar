package chat

import "time"

type Chat struct {
	ID          string     `json:"id"`
	Kind        string     `json:"kind"`
	Title       string     `json:"title"`
	About       string     `json:"about"`
	Origin      *Origin    `json:"origin"`
	PeerID      *string    `json:"peerId"`
	Role        string     `json:"role"`
	Muted       bool       `json:"muted"`
	Seq         int64      `json:"seq"`
	Rev         int64      `json:"rev"`
	ReadSeq     int64      `json:"readSeq"`
	Unread      int64      `json:"unread"`
	MemberCount int        `json:"memberCount"`
	Last        *Message   `json:"last"`
	ArchivedAt  *time.Time `json:"archivedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	ChangedAt   time.Time  `json:"changedAt"`
}

func (c Chat) writable() bool {
	return c.ArchivedAt == nil && (c.Kind != KindChannel || manages(c.Role))
}
