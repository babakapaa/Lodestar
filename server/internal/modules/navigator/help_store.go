package navigator

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const helpScope = "assignment"

type HelpStore struct {
	sealer Sealer
}

func (s HelpStore) Help(ctx context.Context, tx pgx.Tx, userID, scopeKind, scopeID string) (module.HelpSummary, error) {
	rows, err := tx.Query(ctx, `
		SELECT data FROM navigator.help
		WHERE owner_id::text = $1 AND scope_kind = $2 AND scope_id = $3`, userID, scopeKind, scopeID)
	if err != nil {
		return module.HelpSummary{}, fmt.Errorf("load help: %w", err)
	}
	sum, err := pgx.CollectExactlyOneRow(rows, s.scan)
	if errors.Is(err, pgx.ErrNoRows) {
		return emptyHelp(), nil
	}
	return sum, err
}

func (s HelpStore) Record(ctx context.Context, tx pgx.Tx, userID, scopeKind, scopeID string, body messageBody) error {
	sum, err := s.Help(ctx, tx, userID, scopeKind, scopeID)
	if err != nil {
		return err
	}
	sealed, err := s.sealer.seal(countHelp(sum, body))
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO navigator.help (owner_id, scope_kind, scope_id, data) VALUES ($1, $2, $3, $4)
		ON CONFLICT (tenant_id, owner_id, scope_kind, scope_id) DO UPDATE SET data = EXCLUDED.data, changed_at = now()`,
		userID, scopeKind, scopeID, sealed); err != nil {
		return fmt.Errorf("save help: %w", err)
	}
	return nil
}

func (s HelpStore) List(ctx context.Context, tx pgx.Tx, userID string) ([]HelpEntry, error) {
	rows, err := tx.Query(ctx, `
		SELECT scope_kind, scope_id, data FROM navigator.help WHERE owner_id::text = $1 ORDER BY changed_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("list help: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (HelpEntry, error) {
		var e HelpEntry
		var sealed []byte
		if err := row.Scan(&e.ScopeKind, &e.ScopeID, &sealed); err != nil {
			return HelpEntry{}, err
		}
		if err := s.sealer.open(sealed, &e.Summary); err != nil {
			return HelpEntry{}, err
		}
		return e, nil
	})
}

func (s HelpStore) scan(row pgx.CollectableRow) (module.HelpSummary, error) {
	var sealed []byte
	if err := row.Scan(&sealed); err != nil {
		return module.HelpSummary{}, err
	}
	sum := emptyHelp()
	if err := s.sealer.open(sealed, &sum); err != nil {
		return module.HelpSummary{}, err
	}
	return filledHelp(sum), nil
}
