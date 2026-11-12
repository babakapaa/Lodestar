package delivery

import "time"

type quietHours struct {
	on       bool
	from, to int
	location *time.Location
}

func (q quietHours) until(now time.Time) (time.Time, bool) {
	if !q.on || q.from == q.to {
		return time.Time{}, false
	}
	local := now.In(q.location)
	minute := local.Hour()*60 + local.Minute()
	inside := minute >= q.from && minute < q.to
	if q.from > q.to {
		inside = minute >= q.from || minute < q.to
	}
	if !inside {
		return time.Time{}, false
	}
	end := time.Date(local.Year(), local.Month(), local.Day(), q.to/60, q.to%60, 0, 0, q.location)
	if !end.After(local) {
		end = end.AddDate(0, 0, 1)
	}
	return end, true
}
