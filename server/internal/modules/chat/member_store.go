package chat

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/babakapaa/lodestar/server/internal/core/access"
)

const foreignKeyViolation = "23503"

type MemberStore struct{}

func (MemberStore) List(ctx context.Context, tx pgx.Tx, chatID string, managersOnly bool) ([]Member, error) {
	rows, err := tx.Query(ctx, `
		SELECT user_id::text, role, synced, joined_at FROM chat.members
		WHERE chat_id::text = $1 AND (NOT $2 OR role IN ('owner', 'admin'))
		ORDER BY joined_at, user_id`, chatID, managersOnly)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Member, error) {
		var m Member
		err := row.Scan(&m.UserID, &m.Role, &m.Synced, &m.JoinedAt)
		return m, err
	})
}

func (MemberStore) Get(ctx context.Context, tx pgx.Tx, chatID, userID string) (Member, error) {
	var m Member
	err := tx.QueryRow(ctx, `
		SELECT user_id::text, role, synced, joined_at FROM chat.members
		WHERE chat_id::text = $1 AND user_id::text = $2`, chatID, userID).Scan(&m.UserID, &m.Role, &m.Synced, &m.JoinedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Member{}, access.ErrUserNotFound
	}
	if err != nil {
		return Member{}, fmt.Errorf("load member: %w", err)
	}
	return m, nil
}

func (MemberStore) Audience(ctx context.Context, tx pgx.Tx, chatID string) ([]string, error) {
	return userIDs(ctx, tx, `SELECT user_id::text FROM chat.members WHERE chat_id::text = $1`, chatID)
}

func (MemberStore) Recipients(ctx context.Context, tx pgx.Tx, chatID, authorID string) ([]string, error) {
	return userIDs(ctx, tx, `
		SELECT user_id::text FROM chat.members
		WHERE chat_id::text = $1 AND user_id::text <> $2 AND NOT muted AND NOT pending`, chatID, authorID)
}

func (MemberStore) Add(ctx context.Context, tx pgx.Tx, chatID string, ids []string, role string) ([]string, error) {
	added, err := userIDs(ctx, tx, `
		INSERT INTO chat.members (chat_id, user_id, role, read_seq)
		SELECT c.id, u, $3, c.seq FROM chat.chats c, unnest($2::uuid[]) AS u
		WHERE c.id::text = $1
		ON CONFLICT DO NOTHING RETURNING user_id::text`, chatID, ids, role)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolation {
		return nil, access.ErrUserNotFound
	}
	return added, err
}

func (MemberStore) Sync(ctx context.Context, tx pgx.Tx, chatID string, leads, members []string) ([]string, error) {
	ids := append([]string{}, leads...)
	roles := make([]string, 0, len(leads)+len(members))
	for range leads {
		roles = append(roles, RoleAdmin)
	}
	for _, id := range members {
		if !slices.Contains(leads, id) {
			ids, roles = append(ids, id), append(roles, RoleMember)
		}
	}
	changed, err := userIDs(ctx, tx, `
		INSERT INTO chat.members (chat_id, user_id, role, synced, read_seq)
		SELECT c.id, x.u, x.r, true, c.seq FROM chat.chats c, unnest($2::uuid[], $3::text[]) AS x (u, r)
		WHERE c.id::text = $1
		ON CONFLICT (chat_id, user_id) DO UPDATE SET role = EXCLUDED.role, synced = true
		WHERE chat.members.role <> 'owner' AND (chat.members.role <> EXCLUDED.role OR NOT chat.members.synced)
		RETURNING user_id::text`, chatID, ids, roles)
	if err != nil {
		return nil, err
	}
	removed, err := userIDs(ctx, tx, `
		DELETE FROM chat.members WHERE chat_id::text = $1 AND synced AND NOT (user_id::text = ANY($2))
		RETURNING user_id::text`, chatID, ids)
	return append(changed, removed...), err
}

func (MemberStore) Remove(ctx context.Context, tx pgx.Tx, chatID, userID string) error {
	_, err := tx.Exec(ctx, `DELETE FROM chat.members WHERE chat_id::text = $1 AND user_id::text = $2`, chatID, userID)
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}

func (MemberStore) Read(ctx context.Context, tx pgx.Tx, chatID, userID string, seq int64) error {
	_, err := tx.Exec(ctx, `
		UPDATE chat.members m SET read_seq = least($3, c.seq)
		FROM chat.chats c
		WHERE c.id = m.chat_id AND m.chat_id::text = $1 AND m.user_id::text = $2 AND m.read_seq < least($3, c.seq)`,
		chatID, userID, seq)
	if err != nil {
		return fmt.Errorf("mark read: %w", err)
	}
	return nil
}

func (MemberStore) Mute(ctx context.Context, tx pgx.Tx, chatID, userID string, muted bool) error {
	_, err := tx.Exec(ctx, `UPDATE chat.members SET muted = $3 WHERE chat_id::text = $1 AND user_id::text = $2`, chatID, userID, muted)
	if err != nil {
		return fmt.Errorf("mute chat: %w", err)
	}
	return nil
}

func (MemberStore) PeerRead(ctx context.Context, tx pgx.Tx, chatID, userID string) (*int64, error) {
	var seq int64
	err := tx.QueryRow(ctx, `
		SELECT read_seq FROM chat.members WHERE chat_id::text = $1 AND user_id::text <> $2 LIMIT 1`, chatID, userID).Scan(&seq)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("peer read: %w", err)
	}
	return &seq, nil
}

func userIDs(ctx context.Context, tx pgx.Tx, sql string, args ...any) ([]string, error) {
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("member ids: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (MemberStore) Among(ctx context.Context, tx pgx.Tx, chatID string, ids []string) ([]string, error) {
	return userIDs(ctx, tx, `
		SELECT user_id::text FROM chat.members WHERE chat_id::text = $1 AND user_id::text = ANY($2) AND NOT pending
		ORDER BY user_id`, chatID, ids)
}

func (MemberStore) Pending(ctx context.Context, tx pgx.Tx, chatID, userID string, pending bool) error {
	_, err := tx.Exec(ctx, `UPDATE chat.members SET pending = $3 WHERE chat_id::text = $1 AND user_id::text = $2`, chatID, userID, pending)
	if err != nil {
		return fmt.Errorf("set request: %w", err)
	}
	return nil
}

func (MemberStore) Acquainted(ctx context.Context, tx pgx.Tx, a, b string) (bool, error) {
	var found bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM chat.members x JOIN chat.members y ON y.chat_id = x.chat_id JOIN chat.chats c ON c.id = x.chat_id
			WHERE x.user_id::text = $1 AND y.user_id::text = $2 AND c.kind <> 'direct'
		)`, a, b).Scan(&found)
	if err != nil {
		return false, fmt.Errorf("acquainted: %w", err)
	}
	return found, nil
}
