package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/live"
)

type Notifier struct {
	catalog  Catalog
	store    NotificationStore
	channels ChannelStore
	live     live.Store
}

func NewNotifier(catalog Catalog) Notifier {
	return Notifier{catalog: catalog}
}

func (n Notifier) Send(ctx context.Context, tx pgx.Tx, tenantID string, m Message) error {
	t, err := n.catalog.Lookup(m.Type)
	if err != nil {
		return err
	}
	if len(m.UserIDs) == 0 {
		return nil
	}
	params, err := json.Marshal(paramsOrEmpty(m.Params))
	if err != nil {
		return fmt.Errorf("encode params: %w", err)
	}
	choices, err := n.channels.ForType(ctx, tx, m.UserIDs, m.Type)
	if err != nil {
		return err
	}
	var inApp []string
	for _, userID := range m.UserIDs {
		c, chosen := choices[userID]
		if t.Required || (!chosen && t.InApp) || (chosen && c.InApp) {
			inApp = append(inApp, userID)
		}
	}
	if len(inApp) == 0 {
		return nil
	}
	saved, owners, err := n.store.Insert(ctx, tx, inApp, m, params)
	if err != nil {
		return err
	}
	for i, note := range saved {
		if err := n.live.Publish(ctx, tx, tenantID, []string{owners[i]}, LiveType, note); err != nil {
			return err
		}
	}
	return nil
}

func paramsOrEmpty(p map[string]any) map[string]any {
	if p == nil {
		return map[string]any{}
	}
	return p
}
