package navigator

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type JournalPrivacySection struct {
	journal JournalStore
}

func (JournalPrivacySection) Name() string {
	return "navigator.journal"
}

func (s JournalPrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	return s.journal.List(ctx, tx, userID, "", "")
}
