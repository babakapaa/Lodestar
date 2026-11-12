package files

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const MaxFileSetting = "core.files.max_file_mb"

type Settings interface {
	Value(ctx context.Context, tx pgx.Tx, key string) (json.RawMessage, error)
}

func CoreSettings() []schema.Definition {
	lo, hi := schema.Bounds(1, 1<<20)
	return []schema.Definition{
		{Key: MaxFileSetting, Type: schema.Integer, Min: lo, Max: hi, Default: 1024},
	}
}
