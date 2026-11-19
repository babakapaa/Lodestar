package navigator

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type HelpPrivacySection struct {
	help HelpStore
}

func (HelpPrivacySection) Name() string {
	return "navigator.help"
}

func (s HelpPrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	return s.help.List(ctx, tx, userID)
}
