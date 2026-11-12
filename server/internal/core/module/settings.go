package module

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

type Settings interface {
	Value(ctx context.Context, tx pgx.Tx, key string) (json.RawMessage, error)
	ModuleEnabled(ctx context.Context, tenantID, name string) (bool, error)
}
