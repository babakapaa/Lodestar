package chat

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type MessageService struct {
	k           module.Kernel
	chats       ChatStore
	members     MemberStore
	messages    MessageStore
	attachments AttachmentStore
	blocks      BlockStore
	people      identity.People
	announcer   Announcer
}

func NewMessageService(k module.Kernel) MessageService {
	sealer := Sealer{box: k.Secrets}
	return MessageService{k: k, chats: ChatStore{sealer: sealer}, messages: MessageStore{sealer: sealer}}
}

func (s MessageService) Page(ctx context.Context, p auth.Principal, chatID string, q MessageQuery) (MessagePage, error) {
	if !validID(chatID) {
		return MessagePage{}, ErrChatNotFound
	}
	var page MessagePage
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.chats.ForMember(ctx, tx, p.UserID, chatID)
		if err != nil {
			return err
		}
		if page, err = s.messages.Page(ctx, tx, chatID, q, c.Kind == KindChannel); err != nil {
			return err
		}
		return s.attachments.Fill(ctx, tx, page.Messages)
	})
	return page, err
}

func (s MessageService) Send(ctx context.Context, p auth.Principal, chatID string, in MessageInput) (Message, bool, error) {
	in, err := in.normalized()
	if err != nil {
		return Message{}, false, err
	}
	if !validID(chatID) {
		return Message{}, false, ErrChatNotFound
	}
	var m Message
	var created bool
	err = s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.chats.ForMember(ctx, tx, p.UserID, chatID)
		if err != nil {
			return err
		}
		if old, found, err := s.messages.ByClient(ctx, tx, chatID, p.UserID, in.ClientID); err != nil || found {
			m = old
			if err != nil {
				return err
			}
			return s.attachments.Fill(ctx, tx, []Message{m})
		}
		if err := s.check(ctx, tx, p, c, in); err != nil {
			return err
		}
		mentions, err := s.members.Among(ctx, tx, chatID, slices.DeleteFunc(slices.Clone(in.Mentions), func(id string) bool { return id == p.UserID }))
		if err != nil {
			return err
		}
		seq, rev, err := s.chats.Next(ctx, tx, chatID, true)
		if err != nil {
			return err
		}
		if m, err = s.messages.Insert(ctx, tx, chatID, seq, rev, p.UserID, in, messageBody{Text: in.Text, Mentions: mentions}); err != nil {
			return err
		}
		if err := s.attachments.Attach(ctx, tx, m.ID, in.FileIDs); err != nil {
			return err
		}
		list := []Message{m}
		if err := s.attachments.Fill(ctx, tx, list); err != nil {
			return err
		}
		m = list[0]
		if err := s.members.Read(ctx, tx, chatID, p.UserID, seq); err != nil {
			return err
		}
		if c.Request {
			if err := s.members.Pending(ctx, tx, chatID, p.UserID, false); err != nil {
				return err
			}
		}
		if err := s.announcer.Message(ctx, tx, p.TenantID, m); err != nil {
			return err
		}
		created = true
		return s.notify(ctx, tx, p, c, mentions)
	})
	return m, created, err
}

func (s MessageService) check(ctx context.Context, tx pgx.Tx, p auth.Principal, c Chat, in MessageInput) error {
	if !c.writable() {
		return ErrReadOnly
	}
	if c.Kind == KindDirect && c.PeerID != nil {
		blocked, err := s.blocks.Blocks(ctx, tx, *c.PeerID, p.UserID)
		if err != nil {
			return err
		}
		if blocked {
			return ErrBlocked
		}
	}
	if in.ReplyTo != nil {
		found, err := s.messages.Exists(ctx, tx, c.ID, *in.ReplyTo)
		if err != nil {
			return err
		}
		if !found {
			return ErrInvalidMessage
		}
	}
	return s.attachments.Check(ctx, tx, p.UserID, in.FileIDs)
}

func (s MessageService) Edit(ctx context.Context, p auth.Principal, id string, in EditInput) (Message, error) {
	text, err := normalizedText(in.Text)
	if err != nil {
		return Message{}, err
	}
	return s.change(ctx, p, id, func(tx pgx.Tx, c Chat, m Message) error {
		if m.AuthorID == nil || *m.AuthorID != p.UserID || !c.writable() {
			return access.ErrForbidden
		}
		_, rev, err := s.chats.Next(ctx, tx, c.ID, false)
		if err != nil {
			return err
		}
		return s.messages.Edit(ctx, tx, m.ID, rev, messageBody{Text: text, Mentions: m.Mentions})
	})
}

func (s MessageService) Delete(ctx context.Context, p auth.Principal, id string) (Message, error) {
	return s.change(ctx, p, id, func(tx pgx.Tx, c Chat, m Message) error {
		own := m.AuthorID != nil && *m.AuthorID == p.UserID
		if c.ArchivedAt != nil || (!own && (c.Kind == KindDirect || !manages(c.Role))) {
			return access.ErrForbidden
		}
		_, rev, err := s.chats.Next(ctx, tx, c.ID, false)
		if err != nil {
			return err
		}
		if err := s.attachments.Detach(ctx, tx, []string{m.ID}); err != nil {
			return err
		}
		return s.messages.Delete(ctx, tx, m.ID, rev)
	})
}

func (s MessageService) change(ctx context.Context, p auth.Principal, id string, apply func(tx pgx.Tx, c Chat, m Message) error) (Message, error) {
	if !validID(id) {
		return Message{}, ErrMessageNotFound
	}
	var out Message
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		m, err := s.messages.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		c, err := s.chats.ForMember(ctx, tx, p.UserID, m.ChatID)
		if err != nil || m.Deleted {
			return ErrMessageNotFound
		}
		if err := apply(tx, c, m); err != nil {
			return err
		}
		if out, err = s.messages.ByID(ctx, tx, id); err != nil {
			return err
		}
		list := []Message{out}
		if err := s.attachments.Fill(ctx, tx, list); err != nil {
			return err
		}
		out = list[0]
		return s.announcer.Message(ctx, tx, p.TenantID, out)
	})
	return out, err
}

func (s MessageService) notify(ctx context.Context, tx pgx.Tx, p auth.Principal, c Chat, mentions []string) error {
	recipients, err := s.members.Recipients(ctx, tx, c.ID, p.UserID)
	if err != nil {
		return err
	}
	recipients = slices.DeleteFunc(recipients, func(id string) bool { return slices.Contains(mentions, id) })
	if len(recipients) == 0 && len(mentions) == 0 {
		return nil
	}
	people, err := s.people.ByIDs(ctx, tx, []string{p.UserID})
	if err != nil || len(people) == 0 {
		return err
	}
	title := c.Title
	if c.Kind == KindDirect {
		title = people[0].DisplayName
	}
	params := map[string]any{"chat": title, "author": people[0].DisplayName}
	if len(mentions) > 0 {
		err := s.k.Notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyMention, UserIDs: mentions, Params: params, Link: chatLink(c.ID),
		})
		if err != nil {
			return err
		}
	}
	if len(recipients) == 0 {
		return nil
	}
	return s.k.Notifier.Notify(ctx, tx, p.TenantID, module.Notification{
		Type: NotifyMessage, UserIDs: recipients, Params: params, Link: chatLink(c.ID),
	})
}

func chatLink(id string) string {
	return "/chats/" + id
}
