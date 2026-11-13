package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/babakapaa/lodestar/server/internal/core/access"
)

type BlockStore struct{}

func (BlockStore) Add(ctx context.Context, tx pgx.Tx, userID, blockedID string) error {
	_, err := tx.Exec(ctx, `INSERT INTO chat.blocks (user_id, blocked_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, blockedID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolation {
		return access.ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("block: %w", err)
	}
	return nil
}

func (BlockStore) Remove(ctx context.Context, tx pgx.Tx, userID, blockedID string) error {
	_, err := tx.Exec(ctx, `DELETE FROM chat.blocks WHERE user_id::text = $1 AND blocked_id::text = $2`, userID, blockedID)
	if err != nil {
		return fmt.Errorf("unblock: %w", err)
	}
	return nil
}

func (BlockStore) List(ctx context.Context, tx pgx.Tx, userID string) ([]string, error) {
	return userIDs(ctx, tx, `SELECT blocked_id::text FROM chat.blocks WHERE user_id::text = $1 ORDER BY created_at, blocked_id`, userID)
}

func (BlockStore) Blocks(ctx context.Context, tx pgx.Tx, userID, otherID string) (bool, error) {
	return exists(ctx, tx, `SELECT 1 FROM chat.blocks WHERE user_id::text = $1 AND blocked_id::text = $2`, userID, otherID)
}
