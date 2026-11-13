package league

import "time"

type Season struct {
	ID        string     `json:"id"`
	Number    int        `json:"number"`
	StartedAt time.Time  `json:"startedAt"`
	EndedAt   *time.Time `json:"endedAt"`
}
