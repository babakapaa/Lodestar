package team

import "time"

type Scope struct {
	Kind           string     `json:"kind"`
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Context        string     `json:"context"`
	Link           string     `json:"link"`
	MinSize        int        `json:"minSize"`
	MaxSize        int        `json:"maxSize"`
	SelfFormed     bool       `json:"selfFormed"`
	FreeLeaveUntil *time.Time `json:"freeLeaveUntil"`
	DueAt          *time.Time `json:"dueAt"`
	Archived       bool       `json:"archived"`
}
