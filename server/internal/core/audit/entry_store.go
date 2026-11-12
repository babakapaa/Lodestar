package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const pageSize = 50

type EntryStore struct{}

func (EntryStore) Page(ctx context.Context, tx pgx.Tx, f Filter) ([]Entry, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, actor_id::text, action, target_kind, target_id, details, at FROM core.audit_log
		WHERE ($1 = 0 OR id < $1)
		  AND ($2 = '' OR action = $2 OR action LIKE $2 || '.%')
		  AND ($3 = '' OR actor_id::text = $3)
		ORDER BY id DESC LIMIT $4`, f.Before, f.Action, f.ActorID, pageSize)
	if err != nil {
		return nil, fmt.Errorf("read audit log: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Entry, error) {
		var e Entry
		err := row.Scan(&e.ID, &e.ActorID, &e.Action, &e.TargetKind, &e.TargetID, &e.Details, &e.At)
		return e, err
	})
}

func (EntryStore) Prune(ctx context.Context, tx pgx.Tx, keep time.Duration) error {
	_, err := tx.Exec(ctx, `DELETE FROM core.audit_log WHERE at < now() - $1::interval`, keep)
	return err
}

func (EntryStore) ByActor(ctx context.Context, tx pgx.Tx, actorID string) ([]Entry, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, actor_id::text, action, target_kind, target_id, details, at FROM core.audit_log
		WHERE actor_id::text = $1 ORDER BY id`, actorID)
	if err != nil {
		return nil, fmt.Errorf("export audit log: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Entry, error) {
		var e Entry
		err := row.Scan(&e.ID, &e.ActorID, &e.Action, &e.TargetKind, &e.TargetID, &e.Details, &e.At)
		return e, err
	})
}
