package sandbox

import "time"

type tariff struct {
	discount int
	from     int
	to       int
	loc      *time.Location
}

func (t tariff) night(at time.Time) bool {
	h := at.In(t.loc).Hour()
	switch {
	case t.discount == 0 || t.from == t.to:
		return false
	case t.from < t.to:
		return h >= t.from && h < t.to
	default:
		return h >= t.from || h < t.to
	}
}

func (t tariff) factor(start, end time.Time) float64 {
	if !end.After(start) || t.discount == 0 || t.from == t.to {
		return 1
	}
	var night time.Duration
	for at := start; at.Before(end); {
		next := at.In(t.loc).Truncate(time.Hour).Add(time.Hour)
		if !next.After(at) {
			next = at.Add(time.Hour)
		}
		next = minTime(next, end)
		if t.night(at) {
			night += next.Sub(at)
		}
		at = next
	}
	share := night.Seconds() / end.Sub(start).Seconds()
	return 1 - share*float64(t.discount)/100
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
