package identity

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type People struct{}

func (People) ByIDs(ctx context.Context, tx pgx.Tx, ids []string) ([]Person, error) {
	if len(ids) == 0 {
		return []Person{}, nil
	}
	rows, err := tx.Query(ctx, `
		SELECT id::text, display_name, email FROM core.users
		WHERE id::text = ANY($1) ORDER BY display_name, email`, ids)
	if err != nil {
		return nil, fmt.Errorf("people by ids: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Person, error) {
		var p Person
		err := row.Scan(&p.UserID, &p.DisplayName, &p.Email)
		return p, err
	})
}

func (People) Exists(ctx context.Context, tx pgx.Tx, id string) (bool, error) {
	var found bool
	err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM core.users WHERE id::text = $1)`, id).Scan(&found)
	if err != nil {
		return false, fmt.Errorf("user exists: %w", err)
	}
	return found, nil
}

func (People) IDsByEmail(ctx context.Context, tx pgx.Tx, email string) ([]string, error) {
	rows, err := tx.Query(ctx, `SELECT id::text FROM core.users WHERE lower(email) = lower($1) ORDER BY id`, email)
	if err != nil {
		return nil, fmt.Errorf("people by email: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}
