package study

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

type DefenseSlot struct {
	ID           string            `json:"id"`
	AssignmentID string            `json:"assignmentId"`
	StartsAt     time.Time         `json:"startsAt"`
	Minutes      int               `json:"minutes"`
	Place        string            `json:"place"`
	Capacity     int               `json:"capacity"`
	Taken        int               `json:"taken"`
	Mine         bool              `json:"mine"`
	Students     []identity.Person `json:"students,omitempty"`
}

func (s DefenseSlot) Free() bool {
	return s.Taken < s.Capacity
}
