package module

import "time"

type TeamScope struct {
	Kind           string
	ID             string
	Title          string
	Context        string
	Link           string
	Leads          []string
	Members        []string
	MinSize        int
	MaxSize        int
	SelfFormed     bool
	FreeLeaveUntil *time.Time
	DueAt          *time.Time
	Archived       bool
}
