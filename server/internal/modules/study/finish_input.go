package study

type FinishInput struct {
	Status DefenseStatus `json:"status"`
	Grade  *GradeInput   `json:"grade"`
}

func (in FinishInput) valid() bool {
	return in.Status.finished()
}
