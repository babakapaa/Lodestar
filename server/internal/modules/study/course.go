package study

import "time"

type Course struct {
	ID            string     `json:"id"`
	UnitID        *string    `json:"unitId"`
	Title         string     `json:"title"`
	Code          *string    `json:"code"`
	Description   string     `json:"description"`
	ArchivedAt    *time.Time `json:"archivedAt"`
	NavigatorMode string     `json:"navigatorMode"`
	CreatedAt     time.Time  `json:"createdAt"`
}
