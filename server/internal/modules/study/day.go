package study

import (
	"encoding/json"
	"time"
)

const dayLayout = "2006-01-02"

type Day struct {
	time.Time
}

func ParseDay(value string) (Day, error) {
	t, err := time.Parse(dayLayout, value)
	if err != nil {
		return Day{}, ErrInvalidRange
	}
	return Day{t}, nil
}

func DayOf(t time.Time) Day {
	y, m, d := t.Date()
	return Day{time.Date(y, m, d, 0, 0, 0, 0, time.UTC)}
}

func (d Day) String() string {
	return d.Format(dayLayout)
}

func (d Day) AddDays(n int) Day {
	return Day{d.AddDate(0, 0, n)}
}

func (d Day) Weekday() int {
	w := int(d.Time.Weekday())
	if w == 0 {
		return 7
	}
	return w
}

func (d Day) At(minute int, loc *time.Location) time.Time {
	y, m, dd := d.Date()
	return time.Date(y, m, dd, minute/60, minute%60, 0, 0, loc)
}

func (d Day) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Day) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParseDay(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}
