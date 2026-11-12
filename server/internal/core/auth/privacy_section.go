package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct{}

type exportedSession struct {
	DeviceName string     `json:"deviceName"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsedAt time.Time  `json:"lastUsedAt"`
	RevokedAt  *time.Time `json:"revokedAt"`
}

func (PrivacySection) Name() string {
	return "sessions"
}

func (PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	rows, err := tx.Query(ctx, `
		SELECT device_name, created_at, last_used_at, revoked_at FROM core.sessions
		WHERE user_id::text = $1 ORDER BY created_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export sessions: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[exportedSession])
}
