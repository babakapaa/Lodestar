package study

import "time"

const maxTeamSize = 12

type TeamRules struct {
	MinSize        int        `json:"minSize"`
	MaxSize        int        `json:"maxSize"`
	SelfFormed     bool       `json:"selfFormed"`
	FreeLeaveUntil *time.Time `json:"freeLeaveUntil"`
}

func (r TeamRules) valid() bool {
	return r.MinSize >= 2 && r.MaxSize >= r.MinSize && r.MaxSize <= maxTeamSize
}
