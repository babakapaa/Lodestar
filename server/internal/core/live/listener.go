package live

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const retryDelay = 2 * time.Second

type Listener struct {
	pool    *pgxpool.Pool
	tenants database.TenantTx
	hub     *Hub
	store   Store
}

func NewListener(pool *pgxpool.Pool, hub *Hub) Listener {
	return Listener{pool: pool, tenants: database.NewTenantTx(pool), hub: hub}
}

func (l Listener) Run(ctx context.Context) {
	for ctx.Err() == nil {
		if err := l.listen(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("live listener restarts", "error", err)
			select {
			case <-ctx.Done():
			case <-time.After(retryDelay):
			}
		}
	}
}

func (l Listener) listen(ctx context.Context) error {
	conn, err := l.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, "LISTEN "+notifyChannel); err != nil {
		return err
	}
	for {
		n, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			return err
		}
		l.deliver(ctx, n.Payload)
	}
}

func (l Listener) deliver(ctx context.Context, payload string) {
	s, err := decodeSignal(payload)
	if err != nil {
		slog.Warn("skip live signal", "error", err)
		return
	}
	if !l.hub.watching(s.tenantID, s.userID) {
		return
	}
	var e Event
	err = l.tenants.Run(ctx, s.tenantID, func(tx pgx.Tx) error {
		var err error
		e, err = l.store.ByID(ctx, tx, s.eventID)
		return err
	})
	if err != nil {
		slog.Warn("load live event", "error", err)
		return
	}
	l.hub.dispatch(s.tenantID, s.userID, e)
}
