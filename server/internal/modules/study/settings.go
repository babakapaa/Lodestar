package study

import (
	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	TeachersEditSchedule = "study.teachers_edit_schedule"
	TimeZoneSetting      = "core.time_zone"
)

func (Module) Settings() []schema.Definition {
	return []schema.Definition{
		{Key: TeachersEditSchedule, Type: schema.Bool, Default: true},
	}
}
