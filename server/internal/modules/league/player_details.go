package league

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type PlayerDetails struct {
	Person      identity.Person `json:"person"`
	Standing    Standing        `json:"standing"`
	Entries     []Entry         `json:"entries"`
	Transitions []Transition    `json:"transitions"`
	Disputes    []Dispute       `json:"disputes"`
	Grants      []Grant         `json:"grants"`
}
