package notify

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	pageSize            = 50
	notificationColumns = `id::text, type, params, link, created_at, read_at`
)

type NotificationStore struct{}

func (NotificationStore) Insert(ctx context.Context, tx pgx.Tx, userIDs []string, m Message, params []byte) ([]Notification, []string, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO core.notifications (user_id, type, params, link)
		SELECT u, $2, $3, $4 FROM unnest($1::uuid[]) AS u
		RETURNING user_id::text, `+notificationColumns, userIDs, m.Type, params, m.Link)
	if err != nil {
		return nil, nil, fmt.Errorf("save notifications: %w", err)
	}
	var list []Notification
	var owners []string
	for rows.Next() {
		var owner string
		var n Notification
		if err := rows.Scan(&owner, &n.ID, &n.Type, &n.Params, &n.Link, &n.CreatedAt, &n.ReadAt); err != nil {
			rows.Close()
			return nil, nil, err
		}
		list = append(list, n)
		owners = append(owners, owner)
	}
	return list, owners, rows.Err()
}

func (NotificationStore) Page(ctx context.Context, tx pgx.Tx, userID string, before *time.Time, unread bool) ([]Notification, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+notificationColumns+` FROM core.notifications
		WHERE user_id::text = $1 AND ($2::timestamptz IS NULL OR created_at < $2) AND (NOT $3 OR read_at IS NULL)
		ORDER BY created_at DESC, id LIMIT $4`, userID, before, unread, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Notification, error) {
		var n Notification
		err := row.Scan(&n.ID, &n.Type, &n.Params, &n.Link, &n.CreatedAt, &n.ReadAt)
		return n, err
	})
}

func (NotificationStore) Unread(ctx context.Context, tx pgx.Tx, userID string) (int, error) {
	var n int
	err := tx.QueryRow(ctx, `SELECT count(*) FROM core.notifications WHERE user_id::text = $1 AND read_at IS NULL`, userID).Scan(&n)
	return n, err
}

func (NotificationStore) MarkRead(ctx context.Context, tx pgx.Tx, userID, id string) (bool, error) {
	tag, err := tx.Exec(ctx, `
		UPDATE core.notifications SET read_at = coalesce(read_at, now())
		WHERE id::text = $1 AND user_id::text = $2`, id, userID)
	return tag.RowsAffected() == 1, err
}

func (NotificationStore) MarkAllRead(ctx context.Context, tx pgx.Tx, userID string) error {
	_, err := tx.Exec(ctx, `UPDATE core.notifications SET read_at = now() WHERE user_id::text = $1 AND read_at IS NULL`, userID)
	return err
}
