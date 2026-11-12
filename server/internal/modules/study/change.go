package study

import "time"

type ChangeStatus string

const (
	ChangeCancelled ChangeStatus = "cancelled"
	ChangeChanged   ChangeStatus = "changed"
)

type Change struct {
	SeriesID  string       `json:"seriesId"`
	Date      Day          `json:"date"`
	Status    ChangeStatus `json:"status"`
	StartsAt  *time.Time   `json:"startsAt"`
	Minutes   *int         `json:"minutes"`
	Room      *string      `json:"room"`
	TeacherID *string      `json:"teacherId"`
	Note      string       `json:"note"`
}
