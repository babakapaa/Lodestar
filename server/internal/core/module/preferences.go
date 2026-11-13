package module

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

type Preferences interface {
	Value(ctx context.Context, tx pgx.Tx, userID, key string) (json.RawMessage, error)
}
