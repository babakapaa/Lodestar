package attendance

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct{}

func (PrivacySection) Name() string {
	return "attendance.marks"
}

func (PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	return MarkStore{}.Mine(ctx, tx, userID)
}
