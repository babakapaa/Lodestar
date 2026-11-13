package chat

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ReportsPrivacySection struct {
	reports ReportStore
}

func (ReportsPrivacySection) Name() string {
	return "chat.reports"
}

func (s ReportsPrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	return s.reports.ByReporter(ctx, tx, userID)
}
