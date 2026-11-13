package sandbox

import "time"

type Term struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	StartsOn string `json:"startsOn"`
	EndsOn   string `json:"endsOn"`
}

func (t Term) bounds() (time.Time, time.Time) {
	from, _ := time.Parse(time.DateOnly, t.StartsOn)
	to, _ := time.Parse(time.DateOnly, t.EndsOn)
	return from, to.AddDate(0, 0, 1)
}
