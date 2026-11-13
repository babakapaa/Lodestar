package sandbox

import (
	"math"
	"strings"
	"unicode/utf8"
)

type GPUClassInput struct {
	Name     string  `json:"name"`
	MemoryGB int     `json:"memoryGb"`
	Rate     float64 `json:"rate"`
	Archived bool    `json:"archived"`
}

func (in GPUClassInput) normalized() (GPUClassInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Rate = hours(in.Rate)
	if in.Name == "" || utf8.RuneCountInString(in.Name) > 100 || in.MemoryGB < 1 || in.MemoryGB > 4096 ||
		in.Rate <= 0 || in.Rate > 1000 || math.IsNaN(in.Rate) {
		return in, ErrInvalidClass
	}
	return in, nil
}
