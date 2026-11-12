package audit

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const RetentionSetting = "core.audit_retention_days"

type Settings interface {
	Value(ctx context.Context, tx pgx.Tx, key string) (json.RawMessage, error)
}

func CoreSettings() []schema.Definition {
	lo, hi := schema.Bounds(30, 3650)
	return []schema.Definition{
		{Key: RetentionSetting, Type: schema.Integer, Min: lo, Max: hi, Default: 365},
	}
}
