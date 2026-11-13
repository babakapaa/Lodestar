package module

import "encoding/json"

type TeamDelivery struct {
	TenantID   string
	ScopeID    string
	TeamID     string
	TeamName   string
	ProposalID string
	ProposedBy string
	Work       json.RawMessage
	Votes      []TeamVote
}
