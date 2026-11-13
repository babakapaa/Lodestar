package stats

import "github.com/babakapaa/lodestar/server/internal/core/module"

type UnitStats struct {
	UnitID  string          `json:"unitId"`
	Kind    string          `json:"kind"`
	Name    string          `json:"name"`
	People  *int            `json:"people"`
	Hidden  bool            `json:"hidden"`
	Metrics []module.Metric `json:"metrics"`
}
