package navigator

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type JournalStore struct {
	sealer Sealer
}

func (s JournalStore) Insert(ctx context.Context, tx pgx.Tx, ownerID string, d Dialog, body journalBody) error {
	sealed, err := s.sealer.seal(body)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO navigator.journal (owner_id, scope_kind, scope_id, data) VALUES ($1, $2, $3, $4)`,
		ownerID, d.ScopeKind, d.ScopeID, sealed); err != nil {
		return fmt.Errorf("save journal: %w", err)
	}
	return nil
}

func (s JournalStore) List(ctx context.Context, tx pgx.Tx, ownerID, scopeKind, scopeID string) ([]JournalEntry, error) {
	rows, err := tx.Query(ctx, `
		SELECT id::text, scope_kind, scope_id, data, created_at FROM navigator.journal
		WHERE owner_id::text = $1 AND ($2 = '' OR (scope_kind = $2 AND scope_id = $3))
		ORDER BY created_at, id`, ownerID, scopeKind, scopeID)
	if err != nil {
		return nil, fmt.Errorf("list journal: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (JournalEntry, error) {
		var e JournalEntry
		var sealed []byte
		if err := row.Scan(&e.ID, &e.ScopeKind, &e.ScopeID, &sealed, &e.CreatedAt); err != nil {
			return JournalEntry{}, err
		}
		var body journalBody
		if err := s.sealer.open(sealed, &body); err != nil {
			return JournalEntry{}, err
		}
		e.Query, e.Sources = body.Query, body.Sources
		if e.Sources == nil {
			e.Sources = []JournalSource{}
		}
		return e, nil
	})
}

func (JournalStore) Delete(ctx context.Context, tx pgx.Tx, ownerID, id string) error {
	tag, err := tx.Exec(ctx, `DELETE FROM navigator.journal WHERE id::text = $1 AND owner_id::text = $2`, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete journal entry: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrJournalNotFound
	}
	return nil
}
