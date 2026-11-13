package navigator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type DialogStore struct {
	sealer Sealer
}

const dialogColumns = `id::text, scope_kind, scope_id, title, created_at, changed_at`

func (s DialogStore) Insert(ctx context.Context, tx pgx.Tx, ownerID string, in DialogInput) (Dialog, bool, error) {
	tag, err := tx.Exec(ctx, `
		INSERT INTO navigator.dialogs (id, owner_id, scope_kind, scope_id) VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING`, in.ID, ownerID, in.ScopeKind, in.ScopeID)
	if err != nil {
		return Dialog{}, false, fmt.Errorf("create dialog: %w", err)
	}
	d, err := s.Owned(ctx, tx, ownerID, in.ID, false)
	if errors.Is(err, ErrDialogNotFound) {
		return Dialog{}, false, ErrInvalidDialog
	}
	if err == nil && (d.ScopeKind != in.ScopeKind || d.ScopeID != in.ScopeID) {
		return Dialog{}, false, ErrInvalidDialog
	}
	return d, tag.RowsAffected() == 1, err
}

func (s DialogStore) Owned(ctx context.Context, tx pgx.Tx, ownerID, id string, lock bool) (Dialog, error) {
	query := `SELECT ` + dialogColumns + ` FROM navigator.dialogs WHERE id::text = $1 AND owner_id::text = $2`
	if lock {
		query += ` FOR UPDATE`
	}
	rows, err := tx.Query(ctx, query, id, ownerID)
	if err != nil {
		return Dialog{}, fmt.Errorf("load dialog: %w", err)
	}
	d, err := pgx.CollectExactlyOneRow(rows, s.scan)
	if errors.Is(err, pgx.ErrNoRows) {
		return Dialog{}, ErrDialogNotFound
	}
	return d, err
}

func (s DialogStore) List(ctx context.Context, tx pgx.Tx, ownerID, scopeKind, scopeID string) ([]Dialog, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+dialogColumns+` FROM navigator.dialogs
		WHERE owner_id::text = $1 AND ($2 = '' OR (scope_kind = $2 AND scope_id = $3))
		ORDER BY changed_at DESC, id`, ownerID, scopeKind, scopeID)
	if err != nil {
		return nil, fmt.Errorf("list dialogs: %w", err)
	}
	return pgx.CollectRows(rows, s.scan)
}

func (s DialogStore) Touch(ctx context.Context, tx pgx.Tx, id, title string) (time.Time, error) {
	var sealed []byte
	if title != "" {
		var err error
		if sealed, err = s.sealer.seal(title); err != nil {
			return time.Time{}, err
		}
	}
	var at time.Time
	err := tx.QueryRow(ctx, `
		UPDATE navigator.dialogs SET changed_at = now(), title = coalesce(title, $2)
		WHERE id::text = $1 RETURNING changed_at`, id, sealed).Scan(&at)
	if err != nil {
		return time.Time{}, fmt.Errorf("touch dialog: %w", err)
	}
	return at, nil
}

func (DialogStore) Delete(ctx context.Context, tx pgx.Tx, ownerID, id string) error {
	tag, err := tx.Exec(ctx, `DELETE FROM navigator.dialogs WHERE id::text = $1 AND owner_id::text = $2`, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete dialog: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrDialogNotFound
	}
	return nil
}

func (s DialogStore) scan(row pgx.CollectableRow) (Dialog, error) {
	var d Dialog
	var title []byte
	if err := row.Scan(&d.ID, &d.ScopeKind, &d.ScopeID, &title, &d.CreatedAt, &d.ChangedAt); err != nil {
		return Dialog{}, err
	}
	if title != nil {
		if err := s.sealer.open(title, &d.Title); err != nil {
			return Dialog{}, err
		}
	}
	return d, nil
}
