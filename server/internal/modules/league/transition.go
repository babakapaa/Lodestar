package league

import "time"

type Transition struct {
	ID        string           `json:"id"`
	SeasonID  string           `json:"seasonId"`
	UserID    string           `json:"userId"`
	From      int              `json:"from"`
	To        int              `json:"to"`
	Points    int              `json:"points"`
	Status    TransitionStatus `json:"status"`
	DecidedAt *time.Time       `json:"decidedAt"`
	CreatedAt time.Time        `json:"createdAt"`
}

func (t Transition) up() bool {
	return t.To > t.From
}
