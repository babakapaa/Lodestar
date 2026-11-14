package setup

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type HardwareReport struct {
	CheckedAt *time.Time             `json:"checkedAt"`
	Checks    []module.HardwareCheck `json:"checks"`
}

func (r HardwareReport) failed() int {
	n := 0
	for _, c := range r.Checks {
		if c.Required && c.Status != module.CheckOK {
			n++
		}
	}
	return n
}
