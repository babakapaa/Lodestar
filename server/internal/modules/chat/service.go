package chat

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type Service struct {
	k         module.Kernel
	chats     ChatStore
	members   MemberStore
	people    identity.People
	announcer Announcer
	sync      RosterSync
}

func NewService(k module.Kernel) Service {
	return Service{k: k, chats: ChatStore{sealer: Sealer{box: k.Secrets}}, sync: newRosterSync(k)}
}

func newRosterSync(k module.Kernel) RosterSync {
	return RosterSync{chats: ChatStore{sealer: Sealer{box: k.Secrets}}}
}

func (s Service) List(ctx context.Context, p auth.Principal) ([]Chat, error) {
	var list []Chat
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		if list, err = s.chats.List(ctx, tx, p.UserID); err != nil {
			return err
		}
		return s.titleDirects(ctx, tx, list)
	})
	return list, err
}

func (s Service) Get(ctx context.Context, p auth.Principal, id string) (ChatDetails, error) {
	if !validID(id) {
		return ChatDetails{}, ErrChatNotFound
	}
	var d ChatDetails
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		d, err = s.details(ctx, tx, p, id)
		return err
	})
	return d, err
}

func (s Service) Create(ctx context.Context, p auth.Principal, in ChatInput) (ChatDetails, bool, error) {
	in, err := in.normalized()
	if err != nil {
		return ChatDetails{}, false, err
	}
	if in.Origin != nil {
		return s.createLinked(ctx, p, in)
	}
	var d ChatDetails
	err = s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if in.Kind == KindChannel {
			allowed, err := s.k.Access.AllowedInTenant(ctx, tx, p, string(CreateChannels))
			if err != nil {
				return err
			}
			if !allowed {
				return access.ErrForbidden
			}
		}
		id, _, err := s.chats.Insert(ctx, tx, in.Kind, in.Title, in.About, nil, nil, &p.UserID)
		if err != nil {
			return err
		}
		if _, err := s.members.Add(ctx, tx, id, []string{p.UserID}, RoleOwner); err != nil {
			return err
		}
		invited := slices.DeleteFunc(in.MemberIDs, func(u string) bool { return u == p.UserID })
		if _, err := s.members.Add(ctx, tx, id, invited, RoleMember); err != nil {
			return err
		}
		if err := s.announcer.Changed(ctx, tx, p.TenantID, id); err != nil {
			return err
		}
		d, err = s.details(ctx, tx, p, id)
		return err
	})
	return d, err == nil, err
}

func (s Service) createLinked(ctx context.Context, p auth.Principal, in ChatInput) (ChatDetails, bool, error) {
	source, ok := s.k.Rosters.Find(in.Origin.Kind)
	if !ok {
		return ChatDetails{}, false, ErrOriginNotFound
	}
	owner, _, _ := strings.Cut(source.Kind(), ".")
	enabled, err := s.k.Settings.ModuleEnabled(ctx, p.TenantID, owner)
	if err != nil {
		return ChatDetails{}, false, err
	}
	if !enabled {
		return ChatDetails{}, false, ErrOriginNotFound
	}
	var d ChatDetails
	var created bool
	err = s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, err := source.One(ctx, tx, in.Origin.ID)
		if errors.Is(err, module.ErrRosterNotFound) {
			return ErrOriginNotFound
		}
		if err != nil {
			return err
		}
		if !slices.Contains(r.Leads, p.UserID) {
			return access.ErrForbidden
		}
		if r.Archived {
			return ErrReadOnly
		}
		title := in.Title
		if title == "" {
			title = r.Title
		}
		var id string
		if id, created, err = s.chats.Insert(ctx, tx, in.Kind, title, in.About, in.Origin, nil, &p.UserID); err != nil {
			return err
		}
		if err := s.sync.Sync(ctx, tx, p.TenantID, r); err != nil {
			return err
		}
		d, err = s.details(ctx, tx, p, id)
		return err
	})
	return d, created, err
}

func (s Service) Direct(ctx context.Context, p auth.Principal, in DirectInput) (ChatDetails, bool, error) {
	if !validID(in.UserID) || in.UserID == p.UserID {
		return ChatDetails{}, false, ErrInvalidChat
	}
	pair := []string{p.UserID, in.UserID}
	slices.Sort(pair)
	key := strings.Join(pair, ":")
	var d ChatDetails
	var created bool
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		exists, err := s.people.Exists(ctx, tx, in.UserID)
		if err != nil {
			return err
		}
		if !exists {
			return access.ErrUserNotFound
		}
		var id string
		if id, created, err = s.chats.Insert(ctx, tx, KindDirect, "", "", nil, &key, &p.UserID); err != nil {
			return err
		}
		if created {
			if _, err := s.members.Add(ctx, tx, id, pair, RoleMember); err != nil {
				return err
			}
			if err := s.announcer.Changed(ctx, tx, p.TenantID, id); err != nil {
				return err
			}
		}
		d, err = s.details(ctx, tx, p, id)
		return err
	})
	return d, created, err
}

func (s Service) Update(ctx context.Context, p auth.Principal, id string, in ChatUpdate) (ChatDetails, error) {
	in, err := in.normalized()
	if err != nil {
		return ChatDetails{}, err
	}
	var d ChatDetails
	err = s.managed(ctx, p, id, func(tx pgx.Tx, c Chat) error {
		if err := s.chats.Update(ctx, tx, id, in); err != nil {
			return err
		}
		if err := s.announcer.Changed(ctx, tx, p.TenantID, id); err != nil {
			return err
		}
		d, err = s.details(ctx, tx, p, id)
		return err
	})
	return d, err
}

func (s Service) AddMembers(ctx context.Context, p auth.Principal, id string, in MembersInput) (ChatDetails, error) {
	ids, err := memberIDs(in.UserIDs)
	if err != nil || len(ids) == 0 {
		return ChatDetails{}, ErrInvalidChat
	}
	var d ChatDetails
	err = s.managed(ctx, p, id, func(tx pgx.Tx, c Chat) error {
		added, err := s.members.Add(ctx, tx, id, ids, RoleMember)
		if err != nil {
			return err
		}
		if len(added) > 0 {
			if err := s.announcer.Changed(ctx, tx, p.TenantID, id); err != nil {
				return err
			}
		}
		d, err = s.details(ctx, tx, p, id)
		return err
	})
	return d, err
}

func (s Service) RemoveMember(ctx context.Context, p auth.Principal, id, userID string) error {
	if !validID(id) || !validID(userID) {
		return ErrChatNotFound
	}
	return s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.chats.ForMember(ctx, tx, p.UserID, id)
		if err != nil {
			return err
		}
		if c.Kind == KindDirect {
			return ErrInvalidChat
		}
		target, err := s.members.Get(ctx, tx, id, userID)
		if err != nil {
			return err
		}
		switch {
		case target.Synced:
			return ErrSyncedMember
		case target.Role == RoleOwner:
			return ErrOwnerLeaves
		case userID != p.UserID && !outranks(c.Role, target.Role):
			return access.ErrForbidden
		}
		if err := s.members.Remove(ctx, tx, id, userID); err != nil {
			return err
		}
		return s.announcer.Changed(ctx, tx, p.TenantID, id, userID)
	})
}

func (s Service) Mute(ctx context.Context, p auth.Principal, id string, in MuteInput) (Chat, error) {
	return s.own(ctx, p, id, func(tx pgx.Tx, _ Chat) error {
		return s.members.Mute(ctx, tx, id, p.UserID, in.Muted)
	})
}

func (s Service) Read(ctx context.Context, p auth.Principal, id string, in ReadInput) (Chat, error) {
	if in.Seq < 0 {
		return Chat{}, ErrInvalidMessage
	}
	return s.own(ctx, p, id, func(tx pgx.Tx, before Chat) error {
		if in.Seq <= before.ReadSeq {
			return nil
		}
		if err := s.members.Read(ctx, tx, id, p.UserID, in.Seq); err != nil {
			return err
		}
		return s.announcer.Read(ctx, tx, p.TenantID, before, p.UserID, min(in.Seq, before.Seq))
	})
}

func (s Service) own(ctx context.Context, p auth.Principal, id string, change func(tx pgx.Tx, c Chat) error) (Chat, error) {
	if !validID(id) {
		return Chat{}, ErrChatNotFound
	}
	var c Chat
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		before, err := s.chats.ForMember(ctx, tx, p.UserID, id)
		if err != nil {
			return err
		}
		if err := change(tx, before); err != nil {
			return err
		}
		if c, err = s.chats.ForMember(ctx, tx, p.UserID, id); err != nil {
			return err
		}
		list := []Chat{c}
		err = s.titleDirects(ctx, tx, list)
		c = list[0]
		return err
	})
	return c, err
}

func (s Service) managed(ctx context.Context, p auth.Principal, id string, change func(tx pgx.Tx, c Chat) error) error {
	if !validID(id) {
		return ErrChatNotFound
	}
	return s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.chats.ForMember(ctx, tx, p.UserID, id)
		if err != nil {
			return err
		}
		if c.Kind == KindDirect {
			return ErrInvalidChat
		}
		if !manages(c.Role) {
			return access.ErrForbidden
		}
		if c.ArchivedAt != nil {
			return ErrReadOnly
		}
		return change(tx, c)
	})
}

func (s Service) details(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (ChatDetails, error) {
	c, err := s.chats.ForMember(ctx, tx, p.UserID, id)
	if err != nil {
		return ChatDetails{}, err
	}
	list := []Chat{c}
	if err := s.titleDirects(ctx, tx, list); err != nil {
		return ChatDetails{}, err
	}
	d := ChatDetails{Chat: list[0]}
	if d.Members, err = s.members.List(ctx, tx, id, c.Kind == KindChannel && !manages(c.Role)); err != nil {
		return ChatDetails{}, err
	}
	if err := s.name(ctx, tx, d.Members); err != nil {
		return ChatDetails{}, err
	}
	if c.Kind == KindDirect {
		d.PeerReadSeq, err = s.members.PeerRead(ctx, tx, id, p.UserID)
	}
	return d, err
}

func (s Service) name(ctx context.Context, tx pgx.Tx, members []Member) error {
	ids := make([]string, len(members))
	for i, m := range members {
		ids[i] = m.UserID
	}
	names, err := s.names(ctx, tx, ids)
	if err != nil {
		return err
	}
	for i := range members {
		members[i].DisplayName = names[members[i].UserID]
	}
	return nil
}

func (s Service) titleDirects(ctx context.Context, tx pgx.Tx, list []Chat) error {
	var ids []string
	for _, c := range list {
		if c.PeerID != nil {
			ids = append(ids, *c.PeerID)
		}
	}
	names, err := s.names(ctx, tx, ids)
	if err != nil {
		return err
	}
	for i := range list {
		if list[i].PeerID != nil {
			list[i].Title = names[*list[i].PeerID]
		}
	}
	return nil
}

func (s Service) names(ctx context.Context, tx pgx.Tx, ids []string) (map[string]string, error) {
	people, err := s.people.ByIDs(ctx, tx, ids)
	if err != nil {
		return nil, err
	}
	names := make(map[string]string, len(people))
	for _, person := range people {
		names[person.UserID] = person.DisplayName
	}
	return names, nil
}

func outranks(actor, target string) bool {
	switch actor {
	case RoleOwner:
		return target != RoleOwner
	case RoleAdmin:
		return target == RoleMember
	}
	return false
}
