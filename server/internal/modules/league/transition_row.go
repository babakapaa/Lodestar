package league

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type TransitionRow struct {
	Transition
	Person identity.Person `json:"person"`
}
