package study

import "github.com/babakapaa/lodestar/server/internal/core/module"

type courseRosterRow struct {
	roster   module.Roster
	students []string
	units    []string
}
