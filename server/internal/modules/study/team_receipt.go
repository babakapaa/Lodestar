package study

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type TeamReceipt struct {
	TeamID     string           `json:"teamId"`
	TeamName   string           `json:"teamName"`
	ProposalID string           `json:"proposalId"`
	ProposedBy *identity.Person `json:"proposedBy"`
	Votes      []ReceiptVote    `json:"votes"`
}
