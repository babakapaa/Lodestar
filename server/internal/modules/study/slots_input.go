package study

import (
	"strings"
	"time"
)

const maxSlotsAtOnce = 50

type SlotsInput struct {
	StartsAt time.Time `json:"startsAt"`
	Minutes  int       `json:"minutes"`
	Place    string    `json:"place"`
	Capacity int       `json:"capacity"`
	Count    int       `json:"count"`
}

func (in SlotsInput) normalized() (SlotsInput, error) {
	in.Place = strings.TrimSpace(in.Place)
	if in.Capacity == 0 {
		in.Capacity = 1
	}
	if in.Count == 0 {
		in.Count = 1
	}
	if in.StartsAt.IsZero() || in.Minutes < 5 || in.Minutes > 240 || len(in.Place) > 500 {
		return in, ErrInvalidSlot
	}
	if in.Capacity < 1 || in.Capacity > 50 || in.Count < 1 || in.Count > maxSlotsAtOnce {
		return in, ErrInvalidSlot
	}
	return in, nil
}

func (in SlotsInput) starts() []time.Time {
	out := make([]time.Time, in.Count)
	for i := range out {
		out[i] = in.StartsAt.Add(time.Duration(i*in.Minutes) * time.Minute)
	}
	return out
}
