package sandbox

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type LedgerStore struct{}

func (LedgerStore) Lock(ctx context.Context, tx pgx.Tx, termID, ownerKind, ownerID string) error {
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('sandbox.ledger:' || $1 || ':' || $2 || ':' || $3))`,
		termID, ownerKind, ownerID); err != nil {
		return fmt.Errorf("lock ledger: %w", err)
	}
	return nil
}

func (LedgerStore) Add(ctx context.Context, tx pgx.Tx, termID, ownerKind, ownerID, kind string, h float64, title, actorID string) error {
	if _, err := tx.Exec(ctx, `
		INSERT INTO sandbox.ledger (term_id, owner_kind, owner_id, kind, hours, title, actor_id)
		VALUES ($1::uuid, $2, $3, $4, $5, $6, nullif($7, '')::uuid)`,
		termID, ownerKind, ownerID, kind, h, title, actorID); err != nil {
		return fmt.Errorf("add ledger entry: %w", err)
	}
	return nil
}

func (LedgerStore) Charge(ctx context.Context, tx pgx.Tx, r Run, h float64) error {
	ownerKind, ownerID := r.owner()
	if _, err := tx.Exec(ctx, `
		INSERT INTO sandbox.ledger (term_id, owner_kind, owner_id, kind, hours, title, actor_id, run_id)
		VALUES ($1::uuid, $2, $3, 'run', $4, $5, $6::uuid, $7::uuid) ON CONFLICT (run_id) WHERE run_id IS NOT NULL DO NOTHING`,
		r.TermID, ownerKind, ownerID, -h, r.Title, r.UserID, r.ID); err != nil {
		return fmt.Errorf("charge run: %w", err)
	}
	return nil
}

func (LedgerStore) Lines(ctx context.Context, tx pgx.Tx, termID, ownerKind, ownerID string) ([]QuotaLine, error) {
	rows, err := tx.Query(ctx, `
		SELECT CASE kind WHEN 'topup' THEN $4 WHEN 'run' THEN $6 ELSE $5 END, title, hours::float8 FROM sandbox.ledger
		WHERE term_id::text = $1 AND owner_kind = $2 AND owner_id = $3 ORDER BY created_at, id`,
		termID, ownerKind, ownerID, SourceTopup, SourceTransfer, SourceRun)
	if err != nil {
		return nil, fmt.Errorf("ledger lines: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (QuotaLine, error) {
		var l QuotaLine
		err := row.Scan(&l.Source, &l.Title, &l.Hours)
		return l, err
	})
}
