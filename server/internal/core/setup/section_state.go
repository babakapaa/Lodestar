package setup

import "time"

type SectionState struct {
	Key        string     `json:"key"`
	Required   bool       `json:"required"`
	Status     string     `json:"status"`
	AssigneeID *string    `json:"assigneeId"`
	DueOn      *string    `json:"dueOn"`
	Note       string     `json:"note"`
	Settings   []string   `json:"settings"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
