package attendance

type MarkInput struct {
	Status *MarkStatus `json:"status"`
}

func (in MarkInput) normalized() (MarkInput, error) {
	if in.Status != nil && !in.Status.valid() {
		return in, ErrInvalidMark
	}
	return in, nil
}
