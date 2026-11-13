package league

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type TaskDetails struct {
	Task      Task              `json:"task"`
	League    string            `json:"league"`
	Kind      string            `json:"kind"`
	OriginID  string            `json:"originId"`
	Role      Role              `json:"role"`
	Solutions []SolutionRow     `json:"solutions"`
	People    []identity.Person `json:"people"`
}
