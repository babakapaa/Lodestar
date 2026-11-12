package module

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Notifier interface {
	Notify(ctx context.Context, tx pgx.Tx, tenantID string, n Notification) error
}
