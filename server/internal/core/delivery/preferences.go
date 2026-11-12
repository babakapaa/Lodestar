package delivery

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	LanguagePreference   = "core.language"
	QuietHoursPreference = "core.quiet_hours"
	QuietFromPreference  = "core.quiet_from"
	QuietToPreference    = "core.quiet_to"
	TimeZonePreference   = "core.time_zone"
	LanguageSetting      = "core.default_language"
	TimeZoneSetting      = "core.time_zone"
)

type Preferences interface {
	Value(ctx context.Context, tx pgx.Tx, userID, key string) (json.RawMessage, error)
}

type Settings interface {
	Value(ctx context.Context, tx pgx.Tx, key string) (json.RawMessage, error)
}

type Directory interface {
	Email(ctx context.Context, tx pgx.Tx, userID string) (string, error)
}

func CorePreferences() []schema.Definition {
	lo, hi := schema.Bounds(0, 24*60-1)
	return []schema.Definition{
		{Key: QuietHoursPreference, Type: schema.Bool, Default: false},
		{Key: QuietFromPreference, Type: schema.Integer, Min: lo, Max: hi, Default: 22 * 60},
		{Key: QuietToPreference, Type: schema.Integer, Min: lo, Max: hi, Default: 7 * 60},
		{Key: TimeZonePreference, Type: schema.TimeZone, Default: ""},
	}
}

func CoreSettings() []schema.Definition {
	return []schema.Definition{
		{Key: TimeZoneSetting, Type: schema.TimeZone, Default: "UTC"},
	}
}
