package navigator

import "time"

const (
	RoleStudent   = "student"
	RoleNavigator = "navigator"
)

type Message struct {
	Seq       int        `json:"seq"`
	Role      string     `json:"role"`
	Text      string     `json:"text"`
	Intent    Intent     `json:"intent"`
	Step      Step       `json:"step"`
	LastStep  bool       `json:"lastStep"`
	Support   string     `json:"support"`
	Citations []Citation `json:"citations"`
	CreatedAt time.Time  `json:"createdAt"`
}
