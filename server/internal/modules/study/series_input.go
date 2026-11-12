package study

import "strings"

const maxSeriesDays = 400

type SeriesInput struct {
	UnitID     *string `json:"unitId"`
	TeacherID  *string `json:"teacherId"`
	Kind       string  `json:"kind"`
	Room       string  `json:"room"`
	StartsOn   Day     `json:"startsOn"`
	EndsOn     Day     `json:"endsOn"`
	Weekday    int     `json:"weekday"`
	StartTime  Clock   `json:"startTime"`
	Minutes    int     `json:"minutes"`
	EveryWeeks int     `json:"everyWeeks"`
}

func (in SeriesInput) normalized() (SeriesInput, error) {
	in.Kind = strings.TrimSpace(in.Kind)
	in.Room = strings.TrimSpace(in.Room)
	if in.EveryWeeks == 0 {
		in.EveryWeeks = 1
	}
	if in.Weekday == 0 && !in.StartsOn.IsZero() {
		in.Weekday = in.StartsOn.Weekday()
	}
	if in.UnitID != nil && *in.UnitID == "" {
		in.UnitID = nil
	}
	if in.TeacherID != nil && *in.TeacherID == "" {
		in.TeacherID = nil
	}
	if in.StartsOn.IsZero() || in.EndsOn.IsZero() || in.EndsOn.Before(in.StartsOn.Time) ||
		in.EndsOn.Sub(in.StartsOn.Time).Hours() > maxSeriesDays*24 || in.Weekday < 1 || in.Weekday > 7 ||
		in.Minutes < 1 || in.Minutes > 1440 || in.EveryWeeks < 1 || in.EveryWeeks > 4 {
		return in, ErrInvalidSeries
	}
	return in, nil
}
