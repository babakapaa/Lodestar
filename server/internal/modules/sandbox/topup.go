package sandbox

import "time"

type Topup struct {
	ID          string     `json:"id"`
	TermID      string     `json:"termId"`
	UserID      string     `json:"userId"`
	DisplayName string     `json:"displayName"`
	Email       string     `json:"email"`
	WorkKind    string     `json:"workKind"`
	WorkID      string     `json:"workId"`
	Title       string     `json:"title"`
	Hours       float64    `json:"hours"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	Granted     *float64   `json:"granted"`
	CreatedAt   time.Time  `json:"createdAt"`
	DecidedAt   *time.Time `json:"decidedAt"`
	Mine        bool       `json:"mine"`
}

func (t Topup) name() string {
	if t.DisplayName != "" {
		return t.DisplayName
	}
	return t.Email
}
