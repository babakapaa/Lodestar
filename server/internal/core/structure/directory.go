package structure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Directory struct{}

func (Directory) UsersUnder(ctx context.Context, tx pgx.Tx, unitIDs []string) ([]string, error) {
	if len(unitIDs) == 0 {
		return nil, nil
	}
	rows, err := tx.Query(ctx, `
		WITH RECURSIVE inside AS (
			SELECT id FROM core.units WHERE id::text = ANY($1) AND archived_at IS NULL
			UNION
			SELECT u.id FROM core.units u JOIN inside ON u.parent_id = inside.id WHERE u.archived_at IS NULL
		)
		SELECT DISTINCT m.user_id::text FROM core.unit_members m JOIN inside ON inside.id = m.unit_id`, unitIDs)
	if err != nil {
		return nil, fmt.Errorf("users under units: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (Directory) UnitsAbove(ctx context.Context, tx pgx.Tx, userID string) ([]string, error) {
	rows, err := tx.Query(ctx, `
		WITH RECURSIVE above AS (
			SELECT u.id, u.parent_id FROM core.unit_members m JOIN core.units u ON u.id = m.unit_id
			WHERE m.user_id::text = $1 AND u.archived_at IS NULL
			UNION
			SELECT u.id, u.parent_id FROM core.units u JOIN above ON above.parent_id = u.id
		)
		SELECT DISTINCT id::text FROM above`, userID)
	if err != nil {
		return nil, fmt.Errorf("units above user: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}
