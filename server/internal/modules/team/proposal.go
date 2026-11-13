package team

import (
	"encoding/json"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const silence = 24 * time.Hour

type Proposal struct {
	ID         string           `json:"id"`
	TeamID     string           `json:"teamId"`
	ProposedBy *identity.Person `json:"proposedBy"`
	Work       json.RawMessage  `json:"work"`
	Note       string           `json:"note"`
	Status     ProposalStatus   `json:"status"`
	SilentAt   time.Time        `json:"silentAt"`
	DueAt      *time.Time       `json:"dueAt"`
	CreatedAt  time.Time        `json:"createdAt"`
	SentAt     *time.Time       `json:"sentAt"`
	Votes      []Vote           `json:"votes"`
}

func (p Proposal) objected() bool {
	for _, v := range p.Votes {
		if v.Decision == module.VoteObjected {
			return true
		}
	}
	return false
}

func (p Proposal) agreed() bool {
	for _, v := range p.Votes {
		if v.Decision != module.VoteAgreed {
			return false
		}
	}
	return true
}

func (p Proposal) ready(now time.Time) bool {
	if p.Status != ProposalOpen {
		return false
	}
	if p.agreed() {
		return true
	}
	if p.objected() {
		return p.DueAt != nil && !now.Before(*p.DueAt)
	}
	return !now.Before(p.SilentAt)
}

func (p Proposal) delivery(tenantID string, t Team) module.TeamDelivery {
	d := module.TeamDelivery{
		TenantID: tenantID, ScopeID: t.ScopeID, TeamID: t.ID, TeamName: t.Name, ProposalID: p.ID, Work: p.Work,
		Votes: make([]module.TeamVote, 0, len(p.Votes)),
	}
	if p.ProposedBy != nil {
		d.ProposedBy = p.ProposedBy.UserID
	}
	for _, v := range p.Votes {
		decision := v.Decision
		if decision == VoteWaiting {
			decision = module.VoteSilent
		}
		d.Votes = append(d.Votes, module.TeamVote{UserID: v.UserID, Decision: decision, Comment: v.Comment})
	}
	return d
}

func silentAt(now time.Time, due *time.Time) time.Time {
	at := now.Add(silence)
	if due != nil && due.Before(at) {
		return *due
	}
	return at
}
