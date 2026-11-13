package league

type TaskDetails struct {
	Task      Task          `json:"task"`
	League    string        `json:"league"`
	Role      Role          `json:"role"`
	Solutions []SolutionRow `json:"solutions"`
}
