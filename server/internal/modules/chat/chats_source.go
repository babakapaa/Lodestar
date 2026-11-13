package chat

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/search"
)

type ChatsSource struct{}

func (ChatsSource) Kind() string {
	return "chat.chats"
}

func (ChatsSource) Search(ctx context.Context, tx pgx.Tx, p auth.Principal, q search.Query, limit int) ([]search.Hit, error) {
	rows, err := tx.Query(ctx, `
		SELECT c.id::text, c.title, c.kind, word_similarity($1, core.search_norm(c.title)) AS score
		FROM chat.chats c JOIN chat.members m ON m.chat_id = c.id AND m.user_id::text = $2
		WHERE c.kind <> 'direct' AND $1 <% core.search_norm(c.title)
		ORDER BY score DESC, c.title LIMIT $3`, q.Text, p.UserID, limit)
	if err != nil {
		return nil, fmt.Errorf("search chats: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (search.Hit, error) {
		var h search.Hit
		err := row.Scan(&h.ID, &h.Title, &h.Subtitle, &h.Score)
		h.Link = chatLink(h.ID)
		return h, err
	})
}
