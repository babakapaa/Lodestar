package chat

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type FileLink struct{}

func (FileLink) Readable(ctx context.Context, tx pgx.Tx, p auth.Principal, fileID string) (bool, error) {
	return exists(ctx, tx, `
		SELECT 1 FROM chat.message_files f
		JOIN chat.messages m ON m.id = f.message_id
		JOIN chat.members u ON u.chat_id = m.chat_id AND u.user_id::text = $2
		WHERE f.file_id::text = $1`, fileID, p.UserID)
}

func (FileLink) InUse(ctx context.Context, tx pgx.Tx, fileID string) (bool, error) {
	return exists(ctx, tx, `SELECT 1 FROM chat.message_files WHERE file_id::text = $1`, fileID)
}

func exists(ctx context.Context, tx pgx.Tx, sql string, args ...any) (bool, error) {
	var found bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (`+sql+`)`, args...).Scan(&found); err != nil {
		return false, fmt.Errorf("chat file link: %w", err)
	}
	return found, nil
}
