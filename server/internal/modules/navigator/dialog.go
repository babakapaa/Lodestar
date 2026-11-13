package navigator

import "time"

type Dialog struct {
	ID        string    `json:"id"`
	ScopeKind string    `json:"scopeKind"`
	ScopeID   string    `json:"scopeId"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	ChangedAt time.Time `json:"changedAt"`
}
