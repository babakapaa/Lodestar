package league

import "time"

type Task struct {
	ID         string     `json:"id"`
	LeagueID   string     `json:"leagueId"`
	Title      string     `json:"title"`
	About      string     `json:"about"`
	Partner    string     `json:"partner"`
	Points     int        `json:"points"`
	MinSize    int        `json:"minSize"`
	MaxSize    int        `json:"maxSize"`
	SelfFormed bool       `json:"selfFormed"`
	DueAt      *time.Time `json:"dueAt"`
	ClosedAt   *time.Time `json:"closedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}
