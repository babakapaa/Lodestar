package attendance

import "time"

type LessonSummary struct {
	SeriesID  string    `json:"seriesId"`
	Date      string    `json:"date"`
	StartsAt  time.Time `json:"startsAt"`
	Room      string    `json:"room"`
	Cancelled bool      `json:"cancelled"`
	SessionID *string   `json:"sessionId"`
	Open      bool      `json:"open"`
	Counts    Counts    `json:"counts"`
}
