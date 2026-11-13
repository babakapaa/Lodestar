package chat

import "github.com/babakapaa/lodestar/server/internal/core/files"

func (s Sealer) message(r messageRow) (*Message, error) {
	if r.id == nil {
		return nil, nil
	}
	m := &Message{
		ID: *r.id, ChatID: *r.chatID, Seq: *r.seq, Rev: *r.rev, AuthorID: r.authorID, ClientID: r.clientID,
		ReplyTo: r.replyTo, Mentions: []string{}, Files: []files.File{}, Views: r.views, Deleted: r.data == nil, EditedAt: r.editedAt, CreatedAt: *r.createdAt,
	}
	if r.data != nil {
		var body messageBody
		if err := s.open(r.data, &body); err != nil {
			return nil, err
		}
		m.Text = body.Text
		if body.Mentions != nil {
			m.Mentions = body.Mentions
		}
	}
	return m, nil
}
