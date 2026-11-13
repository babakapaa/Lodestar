package team

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type ScopeView struct {
	Scope    Scope             `json:"scope"`
	Role     ScopeRole         `json:"role"`
	MyTeam   *Team             `json:"myTeam"`
	Solo     bool              `json:"solo"`
	Seeking  bool              `json:"seeking"`
	Seekers  []Seeker          `json:"seekers"`
	Teams    []Team            `json:"teams"`
	Unplaced []identity.Person `json:"unplaced"`
	Solos    []identity.Person `json:"solos"`
	Requests []LeaveRequest    `json:"requests"`
}
