package study

import (
	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	TeachersEditSchedule = "study.teachers_edit_schedule"
	TimeZoneSetting      = "core.time_zone"
	DefenseRecording     = "study.defense_recording"
)

func (Module) Settings() []schema.Definition {
	return []schema.Definition{
		{Key: TeachersEditSchedule, Type: schema.Bool, Default: true},
		{Key: DefenseRecording, Type: schema.Bool, Default: false},
	}
}
