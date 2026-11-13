package navigator

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

const (
	indexEvery = 15 * time.Second
	indexBatch = 16
)

type Indexer struct {
	pool     *pgxpool.Pool
	tenants  database.TenantTx
	settings module.Settings
	embedder llm.Embedder
	passages PassageStore
	every    time.Duration
}

func NewIndexer(k module.Kernel) Indexer {
	return Indexer{pool: k.Pool, tenants: k.Tenants, settings: k.Settings, embedder: k.Embedder, every: indexEvery}
}

func (x Indexer) Run(ctx context.Context) {
	ticker := time.NewTicker(x.every)
	defer ticker.Stop()
	for {
		if err := x.IndexAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("index navigator materials", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (x Indexer) IndexAll(ctx context.Context) error {
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
		if err != nil || !enabled {
			continue
		}
		if err := x.Index(ctx, tenantID); err != nil {
			return err
		}
	}
	return nil
}

func (x Indexer) Index(ctx context.Context, tenantID string) error {
	model := x.embedder.Name()
	for {
		var pending []passage
		err := x.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			var err error
			pending, err = x.passages.Pending(ctx, tx, model, indexBatch)
			return err
		})
		if err != nil || len(pending) == 0 {
			return err
		}
		texts := make([]string, len(pending))
		for i, p := range pending {
			texts[i] = p.Text
		}
		vectors, err := x.embedder.Embed(ctx, texts)
		if err != nil {
			return err
		}
		err = x.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			for i, p := range pending {
				if err := x.passages.Embed(ctx, tx, model, p, vectors[i]); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil || len(pending) < indexBatch {
			return err
		}
	}
}
