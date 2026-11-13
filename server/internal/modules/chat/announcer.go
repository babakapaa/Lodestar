package chat

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/live"
)

type Announcer struct {
	members MemberStore
	live    live.Store
}

func (a Announcer) Changed(ctx context.Context, tx pgx.Tx, tenantID, chatID string, also ...string) error {
	audience, err := a.members.Audience(ctx, tx, chatID)
	if err != nil {
		return err
	}
	return a.live.Publish(ctx, tx, tenantID, append(audience, also...), ChangedEvent, map[string]any{"chatId": chatID})
}

func (a Announcer) ChangedFor(ctx context.Context, tx pgx.Tx, tenantID, chatID string, userIDs []string) error {
	return a.live.Publish(ctx, tx, tenantID, userIDs, ChangedEvent, map[string]any{"chatId": chatID})
}

func (a Announcer) Message(ctx context.Context, tx pgx.Tx, tenantID string, m Message) error {
	audience, err := a.members.Audience(ctx, tx, m.ChatID)
	if err != nil {
		return err
	}
	return a.live.Publish(ctx, tx, tenantID, audience, MessageEvent, map[string]any{"chatId": m.ChatID, "seq": m.Seq, "rev": m.Rev})
}

func (a Announcer) Read(ctx context.Context, tx pgx.Tx, tenantID string, c Chat, userID string, seq int64) error {
	audience := []string{userID}
	if c.Kind == KindDirect {
		var err error
		if audience, err = a.members.Audience(ctx, tx, c.ID); err != nil {
			return err
		}
	}
	return a.live.Publish(ctx, tx, tenantID, audience, ReadEvent, map[string]any{"chatId": c.ID, "userId": userID, "seq": seq})
}
