package delivery

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Queue struct{}

func (Queue) Enqueue(ctx context.Context, tx pgx.Tx, channel string, userIDs []string, item Item) error {
	if len(userIDs) == 0 {
		return nil
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO core.deliveries (user_id, channel, type, params, link)
		SELECT u, $2, $3, $4, $5 FROM unnest($1::uuid[]) AS u`,
		userIDs, channel, item.Type, []byte(item.Params), item.Link)
	if err != nil {
		return fmt.Errorf("queue %s delivery: %w", channel, err)
	}
	return nil
}
