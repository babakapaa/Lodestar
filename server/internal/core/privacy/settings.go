package privacy

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	SelfDeletionSetting = "core.self_deletion"
	GraceDaysSetting    = "core.deletion_grace_days"
)

type Settings interface {
	Value(ctx context.Context, tx pgx.Tx, key string) (json.RawMessage, error)
}

func CoreSettings() []schema.Definition {
	lo, hi := schema.Bounds(0, 90)
	return []schema.Definition{
		{Key: SelfDeletionSetting, Type: schema.Bool, Default: true},
		{Key: GraceDaysSetting, Type: schema.Integer, Min: lo, Max: hi, Default: 14},
	}
}
