package study

import (
	"strings"
	"time"
)

type ChangeInput struct {
	Status    ChangeStatus `json:"status"`
	StartsAt  *time.Time   `json:"startsAt"`
	Minutes   *int         `json:"minutes"`
	Room      *string      `json:"room"`
	TeacherID *string      `json:"teacherId"`
	Note      string       `json:"note"`
}

func (in ChangeInput) normalized() (ChangeInput, error) {
	in.Note = strings.TrimSpace(in.Note)
	if in.Room != nil {
		room := strings.TrimSpace(*in.Room)
		in.Room = &room
	}
	switch in.Status {
	case ChangeCancelled:
		in.StartsAt, in.Minutes, in.Room, in.TeacherID = nil, nil, nil, nil
	case ChangeChanged:
		if in.StartsAt == nil && in.Minutes == nil && in.Room == nil && in.TeacherID == nil {
			return in, ErrInvalidChange
		}
		if in.Minutes != nil && (*in.Minutes < 1 || *in.Minutes > 1440) {
			return in, ErrInvalidChange
		}
	default:
		return in, ErrInvalidChange
	}
	return in, nil
}
