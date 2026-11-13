package study

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type DefenseRow struct {
	Student  identity.Person `json:"student"`
	Selected bool            `json:"selected"`
	Status   DefenseStatus   `json:"status"`
	Slot     *DefenseSlot    `json:"slot"`
	Grade    *Grade          `json:"grade"`
}
