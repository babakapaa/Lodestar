package navigator

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"
)

type MessageStore struct {
	sealer Sealer
}

const messageColumns = `seq, role, data, created_at`

func (s MessageStore) Insert(ctx context.Context, tx pgx.Tx, dialogID, key, role string, body messageBody) (Message, error) {
	sealed, err := s.sealer.seal(body)
	if err != nil {
		return Message{}, err
	}
	var k *string
	if key != "" {
		k = &key
	}
	rows, err := tx.Query(ctx, `
		INSERT INTO navigator.messages (dialog_id, seq, client_key, role, data)
		VALUES ($1, (SELECT coalesce(max(seq), 0) + 1 FROM navigator.messages WHERE dialog_id = $1), $2, $3, $4)
		RETURNING `+messageColumns, dialogID, k, role, sealed)
	if err != nil {
		return Message{}, fmt.Errorf("save message: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, s.scan)
}

func (s MessageStore) Answered(ctx context.Context, tx pgx.Tx, dialogID, key string) (Answer, bool, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+messageColumns+` FROM navigator.messages
		WHERE dialog_id::text = $1 AND seq IN (
			SELECT seq FROM navigator.messages WHERE dialog_id::text = $1 AND client_key = $2
			UNION SELECT seq + 1 FROM navigator.messages WHERE dialog_id::text = $1 AND client_key = $2)
		ORDER BY seq`, dialogID, key)
	if err != nil {
		return Answer{}, false, fmt.Errorf("find answer: %w", err)
	}
	list, err := pgx.CollectRows(rows, s.scan)
	if err != nil || len(list) != 2 {
		return Answer{}, false, err
	}
	return Answer{Question: list[0], Reply: list[1]}, true, nil
}

func (s MessageStore) List(ctx context.Context, tx pgx.Tx, dialogID string) ([]Message, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+messageColumns+` FROM navigator.messages WHERE dialog_id::text = $1 ORDER BY seq`, dialogID)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	return pgx.CollectRows(rows, s.scan)
}

func (s MessageStore) Memory(ctx context.Context, tx pgx.Tx, ownerID string, d Dialog, limit int) ([]Message, error) {
	rows, err := tx.Query(ctx, `
		SELECT m.seq, m.role, m.data, m.created_at FROM navigator.messages m
		JOIN navigator.dialogs d ON d.id = m.dialog_id
		WHERE d.owner_id::text = $1 AND (d.id::text = $2 OR ($3 <> '' AND d.scope_kind = $3 AND d.scope_id = $4))
		ORDER BY m.created_at DESC, m.seq DESC LIMIT $5`, ownerID, d.ID, d.ScopeKind, d.ScopeID, limit)
	if err != nil {
		return nil, fmt.Errorf("load memory: %w", err)
	}
	list, err := pgx.CollectRows(rows, s.scan)
	slices.Reverse(list)
	return list, err
}

func (s MessageStore) Last(ctx context.Context, tx pgx.Tx, dialogID string) (Message, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+messageColumns+` FROM navigator.messages
		WHERE dialog_id::text = $1 AND role = 'navigator' ORDER BY seq DESC LIMIT 1`, dialogID)
	if err != nil {
		return Message{}, fmt.Errorf("last reply: %w", err)
	}
	m, err := pgx.CollectExactlyOneRow(rows, s.scan)
	if errors.Is(err, pgx.ErrNoRows) {
		return Message{}, ErrInvalidMessage
	}
	return m, err
}

func (s MessageStore) scan(row pgx.CollectableRow) (Message, error) {
	var m Message
	var sealed []byte
	if err := row.Scan(&m.Seq, &m.Role, &sealed, &m.CreatedAt); err != nil {
		return Message{}, err
	}
	var body messageBody
	if err := s.sealer.open(sealed, &body); err != nil {
		return Message{}, err
	}
	m.Text, m.Intent, m.Step, m.Support = body.Text, body.Intent, body.Step, body.Support
	m.LastStep = m.Step.Last()
	return m, nil
}
