package navigator

import "time"

type JournalEntry struct {
	ID        string          `json:"id"`
	ScopeKind string          `json:"scopeKind"`
	ScopeID   string          `json:"scopeId"`
	Query     string          `json:"query"`
	Sources   []JournalSource `json:"sources"`
	CreatedAt time.Time       `json:"createdAt"`
}
