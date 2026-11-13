package league

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type DisputeRow struct {
	Dispute
	Person     identity.Person `json:"person"`
	Entry      *Entry          `json:"entry"`
	Transition *Transition     `json:"transition"`
}
