package team

import "time"

type Team struct {
	ID         string     `json:"id"`
	ScopeKind  string     `json:"scopeKind"`
	ScopeID    string     `json:"scopeId"`
	Name       string     `json:"name"`
	Members    []Member   `json:"members"`
	CreatedAt  time.Time  `json:"createdAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
	inviteCode string
}

func (t Team) has(userID string) bool {
	for _, m := range t.Members {
		if m.UserID == userID {
			return true
		}
	}
	return false
}

func (t Team) ids() []string {
	ids := make([]string, 0, len(t.Members))
	for _, m := range t.Members {
		ids = append(ids, m.UserID)
	}
	return ids
}
