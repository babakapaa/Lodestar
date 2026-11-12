package notes

import "github.com/google/uuid"

type SourceInput struct {
	ID string `json:"id"`
	SourceFields
	FileID *string `json:"fileId"`
}

func (in SourceInput) normalized() (SourceInput, error) {
	fields, err := in.SourceFields.normalized()
	if err != nil {
		return in, err
	}
	in.SourceFields = fields
	if in.FileID != nil {
		if _, err := uuid.Parse(*in.FileID); err != nil {
			return in, ErrInvalidSource
		}
	}
	return in, nil
}
