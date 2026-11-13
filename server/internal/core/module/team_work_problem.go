package module

type TeamWorkProblem struct {
	Status int
	Code   string
}

func (p TeamWorkProblem) Error() string {
	return "team work rejected: " + p.Code
}
