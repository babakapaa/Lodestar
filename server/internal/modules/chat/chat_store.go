package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const chatSelect = `
	SELECT c.id::text, c.kind, c.title, c.about, c.origin_kind, c.origin_id, peer.user_id::text,
	       m.role, m.muted, m.pending, c.seq, c.rev, m.read_seq, greatest(c.seq - m.read_seq, 0),
	       (SELECT count(*) FROM chat.members a WHERE a.chat_id = c.id),
	       c.archived_at, c.created_at, c.changed_at,
	       ` + lastColumns + `
	FROM chat.chats c
	JOIN chat.members m ON m.chat_id = c.id AND m.user_id::text = $1
	LEFT JOIN LATERAL (
		SELECT o.user_id FROM chat.members o WHERE c.kind = 'direct' AND o.chat_id = c.id AND o.user_id <> m.user_id LIMIT 1
	) peer ON true
	LEFT JOIN LATERAL (
		SELECT * FROM chat.messages x WHERE x.chat_id = c.id ORDER BY x.seq DESC LIMIT 1
	) last ON true`

const lastColumns = `last.id::text, last.chat_id::text, last.seq, last.rev, last.author_id::text, last.client_id, last.reply_to, last.data, NULL::int, last.edited_at, last.created_at`

type ChatStore struct {
	sealer Sealer
}

type chatRef struct {
	ID       string
	Kind     string
	Archived bool
}

func (s ChatStore) ForMember(ctx context.Context, tx pgx.Tx, userID, id string) (Chat, error) {
	rows, err := tx.Query(ctx, chatSelect+` WHERE c.id::text = $2`, userID, id)
	if err != nil {
		return Chat{}, fmt.Errorf("load chat: %w", err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, s.scan)
	if errors.Is(err, pgx.ErrNoRows) {
		return Chat{}, ErrChatNotFound
	}
	return c, err
}

func (s ChatStore) List(ctx context.Context, tx pgx.Tx, userID string) ([]Chat, error) {
	rows, err := tx.Query(ctx, chatSelect+` ORDER BY c.changed_at DESC, c.id`, userID)
	if err != nil {
		return nil, fmt.Errorf("list chats: %w", err)
	}
	return pgx.CollectRows(rows, s.scan)
}

func (ChatStore) Insert(ctx context.Context, tx pgx.Tx, kind, title, about string, origin *Origin, directKey, createdBy *string) (string, bool, error) {
	var originKind, originID *string
	if origin != nil {
		originKind, originID = &origin.Kind, &origin.ID
	}
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO chat.chats (kind, title, about, origin_kind, origin_id, direct_key, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT DO NOTHING RETURNING id::text`, kind, title, about, originKind, originID, directKey, createdBy).Scan(&id)
	if err == nil {
		return id, true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", false, fmt.Errorf("create chat: %w", err)
	}
	err = tx.QueryRow(ctx, `
		SELECT id::text FROM chat.chats
		WHERE ($1::text IS NOT NULL AND direct_key = $1) OR (origin_kind = $2 AND origin_id = $3 AND kind = $4)`,
		directKey, originKind, originID, kind).Scan(&id)
	if err != nil {
		return "", false, fmt.Errorf("find existing chat: %w", err)
	}
	return id, false, nil
}

func (ChatStore) ByOrigin(ctx context.Context, tx pgx.Tx, origin Origin) ([]chatRef, error) {
	rows, err := tx.Query(ctx, `
		SELECT id::text, kind, archived_at IS NOT NULL FROM chat.chats
		WHERE origin_kind = $1 AND origin_id = $2 ORDER BY kind`, origin.Kind, origin.ID)
	if err != nil {
		return nil, fmt.Errorf("chats by origin: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (chatRef, error) {
		var r chatRef
		err := row.Scan(&r.ID, &r.Kind, &r.Archived)
		return r, err
	})
}

func (ChatStore) Update(ctx context.Context, tx pgx.Tx, id string, in ChatUpdate) error {
	_, err := tx.Exec(ctx, `UPDATE chat.chats SET title = $2, about = $3 WHERE id::text = $1`, id, in.Title, in.About)
	if err != nil {
		return fmt.Errorf("update chat: %w", err)
	}
	return nil
}

func (ChatStore) Archive(ctx context.Context, tx pgx.Tx, id string, archived bool) error {
	_, err := tx.Exec(ctx, `
		UPDATE chat.chats SET archived_at = CASE WHEN $2 THEN now() END
		WHERE id::text = $1 AND (archived_at IS NULL) = $2`, id, archived)
	if err != nil {
		return fmt.Errorf("archive chat: %w", err)
	}
	return nil
}

func (ChatStore) Next(ctx context.Context, tx pgx.Tx, id string, message bool) (int64, int64, error) {
	var seq, rev int64
	err := tx.QueryRow(ctx, `
		UPDATE chat.chats
		SET seq = seq + CASE WHEN $2 THEN 1 ELSE 0 END, rev = rev + 1,
		    changed_at = CASE WHEN $2 THEN now() ELSE changed_at END
		WHERE id::text = $1 RETURNING seq, rev`, id, message).Scan(&seq, &rev)
	if err != nil {
		return 0, 0, fmt.Errorf("advance chat: %w", err)
	}
	return seq, rev, nil
}

func (s ChatStore) scan(row pgx.CollectableRow) (Chat, error) {
	var r chatRow
	c := &r.chat
	targets := []any{
		&c.ID, &c.Kind, &c.Title, &c.About, &r.originKind, &r.originID, &c.PeerID,
		&c.Role, &c.Muted, &c.Request, &c.Seq, &c.Rev, &c.ReadSeq, &c.Unread, &c.MemberCount,
		&c.ArchivedAt, &c.CreatedAt, &c.ChangedAt,
	}
	if err := row.Scan(append(targets, r.last.targets()...)...); err != nil {
		return Chat{}, err
	}
	if r.originKind != nil {
		c.Origin = &Origin{Kind: *r.originKind, ID: *r.originID}
	}
	last, err := s.sealer.message(r.last)
	if err != nil {
		return Chat{}, err
	}
	c.Last = last
	return *c, nil
}
