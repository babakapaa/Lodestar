package live

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Backlog struct {
	tenants database.TenantTx
	store   Store
}

func NewBacklog(tenants database.TenantTx) Backlog {
	return Backlog{tenants: tenants}
}

func (b Backlog) After(ctx context.Context, p auth.Principal, after int64) ([]Event, bool, error) {
	var events []Event
	err := b.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		events, err = b.store.After(ctx, tx, p.UserID, after)
		return err
	})
	if err != nil {
		return nil, false, err
	}
	if len(events) > replayLimit {
		return events[:replayLimit], true, nil
	}
	return events, false, nil
}

func (b Backlog) Latest(ctx context.Context, p auth.Principal) (int64, error) {
	var id int64
	err := b.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		id, err = b.store.Latest(ctx, tx, p.UserID)
		return err
	})
	return id, err
}
