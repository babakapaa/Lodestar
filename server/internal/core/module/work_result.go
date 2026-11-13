package module

import "time"

type WorkResult struct {
	UserID string
	Score  float64
	Max    float64
	At     time.Time
}

func (r WorkResult) Share() float64 {
	if r.Max <= 0 || r.Score <= 0 {
		return 0
	}
	return r.Score / r.Max
}
