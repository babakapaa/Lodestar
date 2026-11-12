package study

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type StudentWork struct {
	Student  identity.Person `json:"student"`
	Latest   *Submission     `json:"latest"`
	Attempts int             `json:"attempts"`
	Grade    *Grade          `json:"grade"`
}
