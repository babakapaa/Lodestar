package notes

import "github.com/google/uuid"

type CreateInput struct {
	ID   string `json:"id"`
	Meta Meta   `json:"meta"`
}

func (in CreateInput) normalized() (CreateInput, error) {
	if _, err := uuid.Parse(in.ID); err != nil {
		return in, ErrInvalidNote
	}
	meta, err := in.Meta.normalized()
	in.Meta = meta
	return in, err
}
