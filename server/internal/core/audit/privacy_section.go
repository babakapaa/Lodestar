package audit

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct {
	entries EntryStore
}

func (PrivacySection) Name() string {
	return "admin_actions"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	list, err := s.entries.ByActor(ctx, tx, userID)
	if list == nil {
		list = []Entry{}
	}
	return list, err
}
