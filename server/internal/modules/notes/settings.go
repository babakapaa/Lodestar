package notes

import (
	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	QuotaSetting     = "notes.quota_mb"
	TrashDaysSetting = "notes.trash_days"
)

func (Module) Settings() []schema.Definition {
	quotaMin, trashMin, trashMax := 1, 1, 3650
	return []schema.Definition{
		{Key: QuotaSetting, Type: schema.Integer, Default: 10240, Min: &quotaMin},
		{Key: TrashDaysSetting, Type: schema.Integer, Default: 30, Min: &trashMin, Max: &trashMax},
	}
}
