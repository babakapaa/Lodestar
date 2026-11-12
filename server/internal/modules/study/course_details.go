package study

import (
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

type CourseDetails struct {
	Course
	Role      Role              `json:"role"`
	CanManage bool              `json:"canManage"`
	Teachers  []identity.Person `json:"teachers"`
	Units     []structure.Unit  `json:"units"`
}
