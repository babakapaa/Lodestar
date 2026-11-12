package profile

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants     database.TenantTx
	catalog     schema.Catalog
	profiles    ProfileStore
	preferences PreferenceStore
}

func NewService(tenants database.TenantTx, catalog schema.Catalog) Service {
	return Service{tenants: tenants, catalog: catalog}
}

func (s Service) Me(ctx context.Context, p auth.Principal) (Profile, error) {
	var out Profile
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		out, err = s.profiles.Get(ctx, tx, p.UserID)
		return err
	})
	return out, err
}

func (s Service) Of(ctx context.Context, p auth.Principal, userID string) (Profile, error) {
	out, err := s.Me(ctx, auth.Principal{TenantID: p.TenantID, UserID: userID})
	if err != nil || userID == p.UserID {
		return out, err
	}
	public := []Contact{}
	for _, c := range out.Contacts {
		if c.Public {
			public = append(public, c)
		}
	}
	out.Contacts = public
	return out, nil
}

func (s Service) Update(ctx context.Context, p auth.Principal, in ProfileInput) (Profile, error) {
	if err := in.validate(); err != nil {
		return Profile{}, err
	}
	var out Profile
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.profiles.Put(ctx, tx, p.UserID, in); err != nil {
			return err
		}
		var err error
		out, err = s.profiles.Get(ctx, tx, p.UserID)
		return err
	})
	return out, err
}

func (s Service) Value(ctx context.Context, tx pgx.Tx, userID, key string) (json.RawMessage, error) {
	d, err := s.catalog.Lookup(key)
	if err != nil {
		return nil, err
	}
	stored, found, err := s.preferences.Get(ctx, tx, userID, key)
	if err != nil {
		return nil, err
	}
	if found && d.Validate(stored) == nil {
		return stored, nil
	}
	return d.DefaultValue(), nil
}

func (s Service) Catalog() []schema.Definition {
	return s.catalog.All()
}

func (s Service) Preferences(ctx context.Context, p auth.Principal) ([]schema.Value, error) {
	var stored map[string]json.RawMessage
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		stored, err = s.preferences.All(ctx, tx, p.UserID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.catalog.Resolve(stored), nil
}

func (s Service) SetPreference(ctx context.Context, p auth.Principal, key string, value json.RawMessage) error {
	d, err := s.catalog.Lookup(key)
	if err != nil {
		return err
	}
	if err := d.Validate(value); err != nil {
		return err
	}
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		return s.preferences.Put(ctx, tx, p.UserID, key, value)
	})
}

func (s Service) ResetPreference(ctx context.Context, p auth.Principal, key string) error {
	if _, err := s.catalog.Lookup(key); err != nil {
		return err
	}
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		return s.preferences.Delete(ctx, tx, p.UserID, key)
	})
}
