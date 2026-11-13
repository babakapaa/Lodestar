package setup

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

type SettingsWriter interface {
	Catalog() []schema.Definition
	Put(ctx context.Context, tx pgx.Tx, userID, key string, value json.RawMessage) error
	Clear(ctx context.Context, tx pgx.Tx, userID, key string) error
}
