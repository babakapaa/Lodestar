package team

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

type LeaveRequest struct {
	ID        string          `json:"id"`
	TeamID    string          `json:"teamId"`
	TeamName  string          `json:"teamName"`
	Student   identity.Person `json:"student"`
	Reason    string          `json:"reason"`
	Status    LeaveStatus     `json:"status"`
	CreatedAt time.Time       `json:"createdAt"`
}
