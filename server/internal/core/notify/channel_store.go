package notify

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ChannelStore struct{}

func (ChannelStore) ForType(ctx context.Context, tx pgx.Tx, userIDs []string, typ string) (map[string]ChannelsInput, error) {
	rows, err := tx.Query(ctx, `
		SELECT user_id::text, in_app, push, email FROM core.notification_channels
		WHERE type = $2 AND user_id::text = ANY($1)`, userIDs, typ)
	if err != nil {
		return nil, fmt.Errorf("read channels: %w", err)
	}
	choices := map[string]ChannelsInput{}
	var userID string
	var c ChannelsInput
	_, err = pgx.ForEachRow(rows, []any{&userID, &c.InApp, &c.Push, &c.Email}, func() error {
		choices[userID] = c
		return nil
	})
	return choices, err
}

func (ChannelStore) ForUser(ctx context.Context, tx pgx.Tx, userID string) (map[string]ChannelsInput, error) {
	rows, err := tx.Query(ctx, `SELECT type, in_app, push, email FROM core.notification_channels WHERE user_id::text = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("read channels: %w", err)
	}
	choices := map[string]ChannelsInput{}
	var typ string
	var c ChannelsInput
	_, err = pgx.ForEachRow(rows, []any{&typ, &c.InApp, &c.Push, &c.Email}, func() error {
		choices[typ] = c
		return nil
	})
	return choices, err
}

func (ChannelStore) Put(ctx context.Context, tx pgx.Tx, userID, typ string, c ChannelsInput) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO core.notification_channels (user_id, type, in_app, push, email) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, type) DO UPDATE SET in_app = EXCLUDED.in_app, push = EXCLUDED.push, email = EXCLUDED.email`,
		userID, typ, c.InApp, c.Push, c.Email)
	if err != nil {
		return fmt.Errorf("save channels: %w", err)
	}
	return nil
}
