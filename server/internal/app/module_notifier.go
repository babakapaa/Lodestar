package app

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/notify"
)

type moduleNotifier struct {
	notifier notify.Notifier
}

func (m moduleNotifier) Notify(ctx context.Context, tx pgx.Tx, tenantID string, n module.Notification) error {
	return m.notifier.Send(ctx, tx, tenantID, notify.Message{Type: n.Type, UserIDs: n.UserIDs, Params: n.Params, Link: n.Link})
}
