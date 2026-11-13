package chat

import "github.com/babakapaa/lodestar/server/internal/core/schema"

const RetentionSetting = "chat.retention_days"

func (Module) Settings() []schema.Definition {
	minDays, maxDays := 0, 36500
	return []schema.Definition{
		{Key: RetentionSetting, Type: schema.Integer, Default: 0, Min: &minDays, Max: &maxDays},
	}
}
