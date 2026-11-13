package chat

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"
)

const messageColumns = `m.id::text, m.chat_id::text, m.seq, m.rev, m.author_id::text, m.client_id, m.reply_to, m.data`

const viewsColumn = `CASE WHEN $2 THEN (
		SELECT count(*)::int FROM chat.members r
		WHERE r.chat_id = m.chat_id AND r.read_seq >= m.seq AND r.user_id IS DISTINCT FROM m.author_id
	) END`

type MessageStore struct {
	sealer Sealer
}

func (s MessageStore) Insert(ctx context.Context, tx pgx.Tx, chatID string, seq, rev int64, authorID string, in MessageInput) (Message, error) {
	data, err := s.sealer.seal(messageBody{Text: in.Text})
	if err != nil {
		return Message{}, err
	}
	rows, err := tx.Query(ctx, `
		INSERT INTO chat.messages AS m (chat_id, seq, rev, author_id, client_id, reply_to, data)
		VALUES ($1, $3, $4, $5, $6, $7, $8)
		RETURNING `+messageColumns+`, `+viewsColumn+`, m.edited_at, m.created_at`,
		chatID, false, seq, rev, authorID, in.ClientID, in.ReplyTo, data)
	if err != nil {
		return Message{}, fmt.Errorf("insert message: %w", err)
	}
	return s.one(rows)
}

func (s MessageStore) ByClient(ctx context.Context, tx pgx.Tx, chatID, authorID, clientID string) (Message, bool, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+messageColumns+`, `+viewsColumn+`, m.edited_at, m.created_at FROM chat.messages m
		WHERE m.chat_id::text = $1 AND m.author_id::text = $3 AND m.client_id = $4`, chatID, false, authorID, clientID)
	if err != nil {
		return Message{}, false, fmt.Errorf("find message by client id: %w", err)
	}
	m, err := s.one(rows)
	if errors.Is(err, ErrMessageNotFound) {
		return Message{}, false, nil
	}
	return m, err == nil, err
}

func (s MessageStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Message, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+messageColumns+`, `+viewsColumn+`, m.edited_at, m.created_at FROM chat.messages m
		WHERE m.id::text = $1 FOR UPDATE`, id, false)
	if err != nil {
		return Message{}, fmt.Errorf("load message: %w", err)
	}
	return s.one(rows)
}

func (MessageStore) Exists(ctx context.Context, tx pgx.Tx, chatID string, seq int64) (bool, error) {
	var found bool
	err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM chat.messages WHERE chat_id::text = $1 AND seq = $2)`, chatID, seq).Scan(&found)
	if err != nil {
		return false, fmt.Errorf("message exists: %w", err)
	}
	return found, nil
}

func (s MessageStore) Edit(ctx context.Context, tx pgx.Tx, id string, rev int64, text string) error {
	data, err := s.sealer.seal(messageBody{Text: text})
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE chat.messages SET data = $3, rev = $2, edited_at = now() WHERE id::text = $1`, id, rev, data)
	if err != nil {
		return fmt.Errorf("edit message: %w", err)
	}
	return nil
}

func (MessageStore) Delete(ctx context.Context, tx pgx.Tx, id string, rev int64) error {
	_, err := tx.Exec(ctx, `UPDATE chat.messages SET data = NULL, rev = $2, deleted_at = now() WHERE id::text = $1`, id, rev)
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}
	return nil
}

func (s MessageStore) Page(ctx context.Context, tx pgx.Tx, chatID string, q MessageQuery, views bool) (MessagePage, error) {
	query := `SELECT ` + messageColumns + `, ` + viewsColumn + `, m.edited_at, m.created_at FROM chat.messages m WHERE m.chat_id::text = $1`
	args := []any{chatID, views, q.Limit + 1}
	switch {
	case q.Since != nil:
		query += ` AND m.rev > $4 ORDER BY m.rev LIMIT $3`
		args = append(args, *q.Since)
	case q.Before != nil:
		query += ` AND m.seq < $4 ORDER BY m.seq DESC LIMIT $3`
		args = append(args, *q.Before)
	default:
		query += ` ORDER BY m.seq DESC LIMIT $3`
	}
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return MessagePage{}, fmt.Errorf("list messages: %w", err)
	}
	list, err := pgx.CollectRows(rows, s.scan)
	if err != nil {
		return MessagePage{}, err
	}
	page := MessagePage{Messages: list, More: len(list) > q.Limit}
	if page.More {
		page.Messages = list[:q.Limit]
	}
	if q.Since == nil {
		slices.Reverse(page.Messages)
	}
	return page, nil
}

func (s MessageStore) ByAuthor(ctx context.Context, tx pgx.Tx, authorID string) ([]Message, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+messageColumns+`, `+viewsColumn+`, m.edited_at, m.created_at FROM chat.messages m
		WHERE m.author_id::text = $1 AND m.data IS NOT NULL ORDER BY m.created_at, m.id`, authorID, false)
	if err != nil {
		return nil, fmt.Errorf("messages by author: %w", err)
	}
	return pgx.CollectRows(rows, s.scan)
}

func (s MessageStore) one(rows pgx.Rows) (Message, error) {
	m, err := pgx.CollectExactlyOneRow(rows, s.scan)
	if errors.Is(err, pgx.ErrNoRows) {
		return Message{}, ErrMessageNotFound
	}
	return m, err
}

func (s MessageStore) scan(row pgx.CollectableRow) (Message, error) {
	var r messageRow
	if err := row.Scan(r.targets()...); err != nil {
		return Message{}, err
	}
	m, err := s.sealer.message(r)
	if err != nil {
		return Message{}, err
	}
	return *m, nil
}
