package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const purgeEvery = time.Hour

type Purger struct {
	pool      *pgxpool.Pool
	tenants   database.TenantTx
	settings  module.Settings
	announcer Announcer
	every     time.Duration
}

func NewPurger(k module.Kernel) Purger {
	return Purger{pool: k.Pool, tenants: k.Tenants, settings: k.Settings, every: purgeEvery}
}

func (x Purger) Run(ctx context.Context) {
	ticker := time.NewTicker(x.every)
	defer ticker.Stop()
	for {
		if err := x.PurgeAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("purge old chat messages", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (x Purger) PurgeAll(ctx context.Context) error {
	rows, err := x.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	tenants, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range tenants {
		enabled, err := x.settings.ModuleEnabled(ctx, tenantID, Name)
		if err != nil {
			return err
		}
		if enabled {
			if err := x.Purge(ctx, tenantID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (x Purger) Purge(ctx context.Context, tenantID string) error {
	var days int
	var chats []string
	err := x.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		raw, err := x.settings.Value(ctx, tx, RetentionSetting)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(raw, &days); err != nil {
			return fmt.Errorf("decode chat retention: %w", err)
		}
		if days <= 0 {
			return nil
		}
		chats, err = userIDs(ctx, tx, `
			SELECT DISTINCT chat_id::text FROM chat.messages
			WHERE data IS NOT NULL AND created_at < now() - make_interval(days => $1)`, days)
		return err
	})
	if err != nil {
		return err
	}
	for _, chatID := range chats {
		err := x.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			return x.expire(ctx, tx, tenantID, chatID, days)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (x Purger) expire(ctx context.Context, tx pgx.Tx, tenantID, chatID string, days int) error {
	var rev int64
	if err := tx.QueryRow(ctx, `SELECT rev FROM chat.chats WHERE id::text = $1 FOR UPDATE`, chatID).Scan(&rev); err != nil {
		return fmt.Errorf("lock chat: %w", err)
	}
	expired, err := userIDs(ctx, tx, `
		WITH old AS (
			SELECT id, $3 + row_number() OVER (ORDER BY seq) AS rev FROM chat.messages
			WHERE chat_id::text = $1 AND data IS NOT NULL AND created_at < now() - make_interval(days => $2)
		)
		UPDATE chat.messages m SET data = NULL, deleted_at = now(), rev = old.rev
		FROM old WHERE m.id = old.id RETURNING m.id::text`, chatID, days, rev)
	if err != nil || len(expired) == 0 {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE chat.chats SET rev = rev + $2 WHERE id::text = $1`, chatID, len(expired)); err != nil {
		return fmt.Errorf("advance chat: %w", err)
	}
	if err := (AttachmentStore{}).Detach(ctx, tx, expired); err != nil {
		return err
	}
	return x.announcer.Changed(ctx, tx, tenantID, chatID)
}
