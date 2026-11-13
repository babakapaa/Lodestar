package chat

import "time"

type chatRow struct {
	chat       Chat
	originKind *string
	originID   *string
	last       messageRow
}

type messageRow struct {
	id        *string
	chatID    *string
	seq       *int64
	rev       *int64
	authorID  *string
	clientID  *string
	replyTo   *int64
	data      []byte
	views     *int
	editedAt  *time.Time
	createdAt *time.Time
}

func (r *messageRow) targets() []any {
	return []any{&r.id, &r.chatID, &r.seq, &r.rev, &r.authorID, &r.clientID, &r.replyTo, &r.data, &r.views, &r.editedAt, &r.createdAt}
}
