package chat

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/livekit"
)

const (
	MaxCallPeople = 50
	ticketTTL     = 6 * time.Hour
)

type CallService struct {
	k         module.Kernel
	chats     ChatStore
	members   MemberStore
	calls     CallStore
	people    identity.People
	announcer Announcer
	now       func() time.Time
}

func NewCallService(k module.Kernel) CallService {
	return CallService{k: k, chats: ChatStore{sealer: Sealer{box: k.Secrets}}, now: time.Now}
}

func (s CallService) Active(ctx context.Context, p auth.Principal, chatID string) (Call, error) {
	if !validID(chatID) {
		return Call{}, ErrChatNotFound
	}
	var out Call
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.chats.ForMember(ctx, tx, p.UserID, chatID); err != nil {
			return err
		}
		c, found, err := s.calls.Active(ctx, tx, chatID)
		if err != nil {
			return err
		}
		if !found {
			return ErrCallNotFound
		}
		out = c
		return nil
	})
	return out, err
}

func (s CallService) Join(ctx context.Context, p auth.Principal, chatID string) (CallTicket, error) {
	if s.k.Calls == nil {
		return CallTicket{}, ErrCallsUnavailable
	}
	if !validID(chatID) {
		return CallTicket{}, ErrChatNotFound
	}
	var out CallTicket
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.chats.ForMember(ctx, tx, p.UserID, chatID)
		if err != nil {
			return err
		}
		if c.Kind == KindChannel || c.ArchivedAt != nil || c.Request {
			return ErrInvalidChat
		}
		if c.MemberCount > MaxCallPeople {
			return ErrCallFull
		}
		var on bool
		if err := s.setting(ctx, tx, CallsSetting, &on); err != nil {
			return err
		}
		if !on {
			return ErrCallsOff
		}
		call, found, err := s.calls.Active(ctx, tx, chatID)
		if err != nil {
			return err
		}
		if !found {
			id, err := s.calls.Start(ctx, tx, chatID, p.UserID)
			if err != nil {
				return err
			}
			call = Call{ID: id, ChatID: chatID}
			if err := s.ring(ctx, tx, p, c); err != nil {
				return err
			}
		}
		if err := s.calls.Join(ctx, tx, call.ID, p.UserID); err != nil {
			return err
		}
		if call, _, err = s.calls.ByID(ctx, tx, call.ID); err != nil {
			return err
		}
		if len(call.People) > MaxCallPeople {
			return ErrCallFull
		}
		share, err := s.canShare(ctx, tx, c)
		if err != nil {
			return err
		}
		record, err := s.canRecord(ctx, tx, c)
		if err != nil {
			return err
		}
		person, err := s.person(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		room := RoomName(p.TenantID, call.ID)
		token, err := s.k.Calls.Token(p.UserID, person, livekit.Grant{Room: room, ScreenShare: share}, ticketTTL, s.now())
		if err != nil {
			return err
		}
		out = CallTicket{Call: call, URL: s.k.Calls.URL(), Room: room, Token: token, ScreenShare: share, Record: record}
		return s.announce(ctx, tx, p.TenantID, chatID, call.ID, true)
	})
	return out, err
}

func (s CallService) Leave(ctx context.Context, p auth.Principal, callID string) error {
	if !validID(callID) {
		return ErrCallNotFound
	}
	return s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		call, found, err := s.calls.ByID(ctx, tx, callID)
		if err != nil {
			return err
		}
		if !found {
			return ErrCallNotFound
		}
		if _, err := s.chats.ForMember(ctx, tx, p.UserID, call.ChatID); err != nil {
			return ErrCallNotFound
		}
		return s.leave(ctx, tx, p.TenantID, call, p.UserID)
	})
}

func (s CallService) Hook(ctx context.Context, tenantID string, e livekit.Event) error {
	_, callID, ok := strings.Cut(e.Room.Name, "_")
	if !ok || !validID(callID) {
		return nil
	}
	return s.k.Tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		call, found, err := s.calls.ByID(ctx, tx, callID)
		if err != nil || !found {
			return err
		}
		switch e.Event {
		case livekit.ParticipantLeft:
			if !validID(e.Participant.Identity) {
				return nil
			}
			return s.leave(ctx, tx, tenantID, call, e.Participant.Identity)
		case livekit.RoomFinished:
			if err := s.calls.End(ctx, tx, call.ID); err != nil {
				return err
			}
			return s.announce(ctx, tx, tenantID, call.ChatID, call.ID, false)
		}
		return nil
	})
}

func (s CallService) leave(ctx context.Context, tx pgx.Tx, tenantID string, call Call, userID string) error {
	if err := s.calls.Leave(ctx, tx, call.ID, userID); err != nil {
		return err
	}
	people, err := s.calls.People(ctx, tx, call.ID)
	if err != nil {
		return err
	}
	if len(people) == 0 {
		if err := s.calls.End(ctx, tx, call.ID); err != nil {
			return err
		}
	}
	return s.announce(ctx, tx, tenantID, call.ChatID, call.ID, len(people) > 0)
}

func (s CallService) announce(ctx context.Context, tx pgx.Tx, tenantID, chatID, callID string, active bool) error {
	audience, err := s.members.Audience(ctx, tx, chatID)
	if err != nil {
		return err
	}
	return s.announcer.live.Publish(ctx, tx, tenantID, audience, CallEvent, map[string]any{"chatId": chatID, "callId": callID, "active": active})
}

func (s CallService) ring(ctx context.Context, tx pgx.Tx, p auth.Principal, c Chat) error {
	recipients, err := s.members.Recipients(ctx, tx, c.ID, p.UserID)
	if err != nil || len(recipients) == 0 {
		return err
	}
	name, err := s.person(ctx, tx, p.UserID)
	if err != nil {
		return err
	}
	title := c.Title
	if c.Kind == KindDirect {
		title = name
	}
	return s.k.Notifier.Notify(ctx, tx, p.TenantID, module.Notification{
		Type: NotifyCall, UserIDs: recipients, Params: map[string]any{"chat": title, "author": name}, Link: chatLink(c.ID),
	})
}

func (s CallService) canShare(ctx context.Context, tx pgx.Tx, c Chat) (bool, error) {
	var policy string
	if err := s.setting(ctx, tx, ScreenShareSetting, &policy); err != nil {
		return false, err
	}
	return policy == ShareEveryone || c.Kind == KindDirect || manages(c.Role), nil
}

func (s CallService) canRecord(ctx context.Context, tx pgx.Tx, c Chat) (bool, error) {
	if s.k.Calls == nil || !s.k.Calls.Records() {
		return false, nil
	}
	var policy string
	if err := s.setting(ctx, tx, RecordingSetting, &policy); err != nil {
		return false, err
	}
	switch policy {
	case RecordEveryone:
		return true, nil
	case RecordAdmins:
		return c.Kind != KindDirect && manages(c.Role), nil
	}
	return false, nil
}

func (s CallService) person(ctx context.Context, tx pgx.Tx, userID string) (string, error) {
	people, err := s.people.ByIDs(ctx, tx, []string{userID})
	if err != nil || len(people) == 0 {
		return "", err
	}
	return people[0].DisplayName, nil
}

func (s CallService) setting(ctx context.Context, tx pgx.Tx, key string, into any) error {
	raw, err := s.k.Settings.Value(ctx, tx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, into)
}

func RoomName(tenantID, callID string) string {
	return tenantID + "_" + callID
}
