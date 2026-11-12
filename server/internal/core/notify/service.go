package notify

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants  database.TenantTx
	catalog  Catalog
	store    NotificationStore
	channels ChannelStore
}

func NewService(tenants database.TenantTx, catalog Catalog) Service {
	return Service{tenants: tenants, catalog: catalog}
}

func (s Service) List(ctx context.Context, p auth.Principal, before *time.Time, unread bool) ([]Notification, error) {
	var list []Notification
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.store.Page(ctx, tx, p.UserID, before, unread)
		return err
	})
	return list, err
}

func (s Service) Unread(ctx context.Context, p auth.Principal) (int, error) {
	var n int
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		n, err = s.store.Unread(ctx, tx, p.UserID)
		return err
	})
	return n, err
}

func (s Service) MarkRead(ctx context.Context, p auth.Principal, id string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		found, err := s.store.MarkRead(ctx, tx, p.UserID, id)
		if err != nil {
			return err
		}
		if !found {
			return ErrNotFound
		}
		return nil
	})
}

func (s Service) MarkAllRead(ctx context.Context, p auth.Principal) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		return s.store.MarkAllRead(ctx, tx, p.UserID)
	})
}

func (s Service) Types() []TypeInfo {
	list := []TypeInfo{}
	for _, t := range s.catalog.All() {
		list = append(list, TypeInfo{Key: t.Key, Required: t.Required, Names: t.Names, Texts: t.Texts})
	}
	return list
}

func (s Service) Channels(ctx context.Context, p auth.Principal) ([]Channels, error) {
	var choices map[string]ChannelsInput
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		choices, err = s.channels.ForUser(ctx, tx, p.UserID)
		return err
	})
	if err != nil {
		return nil, err
	}
	list := []Channels{}
	for _, t := range s.catalog.All() {
		c := Channels{Type: t.Key, Required: t.Required, InApp: t.InApp, Push: t.Push, Email: t.Email}
		if chosen, ok := choices[t.Key]; ok {
			c.InApp, c.Push, c.Email = chosen.InApp || t.Required, chosen.Push, chosen.Email
		}
		list = append(list, c)
	}
	return list, nil
}

func (s Service) SetChannels(ctx context.Context, p auth.Principal, typ string, in ChannelsInput) error {
	t, err := s.catalog.Lookup(typ)
	if err != nil {
		return err
	}
	if t.Required && !in.InApp {
		return ErrRequiredType
	}
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		return s.channels.Put(ctx, tx, p.UserID, typ, in)
	})
}
