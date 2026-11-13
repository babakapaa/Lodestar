package sandbox

import "time"

type Node struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	ClassID   string     `json:"classId"`
	ClassName string     `json:"className"`
	GPUs      int        `json:"gpus"`
	Busy      int        `json:"busy"`
	SeenAt    *time.Time `json:"seenAt"`
	Online    bool       `json:"online"`
	CreatedAt time.Time  `json:"createdAt"`
}
