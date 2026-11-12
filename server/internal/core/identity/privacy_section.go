package identity

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct{}

func (PrivacySection) Name() string {
	return "account"
}

func (PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	return scanUser(tx.QueryRow(ctx, `SELECT `+userColumns+` FROM core.users WHERE id::text = $1`, userID))
}
