package league

import "time"

type Entry struct {
	ID        string    `json:"id"`
	SeasonID  string    `json:"seasonId"`
	UserID    string    `json:"userId"`
	Kind      EntryKind `json:"kind"`
	WorkKind  *string   `json:"workKind"`
	WorkID    *string   `json:"workId"`
	Title     string    `json:"title"`
	Score     *float64  `json:"score"`
	Max       *float64  `json:"max"`
	Weight    *int      `json:"weight"`
	Points    int       `json:"points"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
