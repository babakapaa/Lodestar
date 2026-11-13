package settings

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants  database.TenantTx
	access   access.Service
	catalog  schema.Catalog
	registry *module.Registry
	settings SettingStore
	modules  ModuleStore
	log      audit.Log
}

func NewService(tenants database.TenantTx, acc access.Service, catalog schema.Catalog, registry *module.Registry) Service {
	return Service{tenants: tenants, access: acc, catalog: catalog, registry: registry}
}

func (s Service) Catalog() []schema.Definition {
	return s.catalog.All()
}

func (s Service) Values(ctx context.Context, p auth.Principal) ([]schema.Value, error) {
	var stored map[string]json.RawMessage
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		stored, err = s.settings.All(ctx, tx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.catalog.Resolve(stored), nil
}

func (s Service) Value(ctx context.Context, tx pgx.Tx, key string) (json.RawMessage, error) {
	d, err := s.catalog.Lookup(key)
	if err != nil {
		return nil, err
	}
	stored, found, err := s.settings.Get(ctx, tx, key)
	if err != nil {
		return nil, err
	}
	if found && d.Validate(stored) == nil {
		return stored, nil
	}
	return d.DefaultValue(), nil
}

func (s Service) Set(ctx context.Context, p auth.Principal, key string, value json.RawMessage) error {
	if err := s.check(key, value); err != nil {
		return err
	}
	return s.managed(ctx, p, func(tx pgx.Tx) error {
		return s.Put(ctx, tx, p.UserID, key, value)
	})
}

func (s Service) Reset(ctx context.Context, p auth.Principal, key string) error {
	if _, err := s.catalog.Lookup(key); err != nil {
		return err
	}
	return s.managed(ctx, p, func(tx pgx.Tx) error {
		return s.Clear(ctx, tx, p.UserID, key)
	})
}

func (s Service) Put(ctx context.Context, tx pgx.Tx, userID, key string, value json.RawMessage) error {
	if err := s.check(key, value); err != nil {
		return err
	}
	before, err := s.Value(ctx, tx, key)
	if err != nil {
		return err
	}
	if err := s.settings.Put(ctx, tx, key, value, userID); err != nil {
		return err
	}
	return s.log.Record(ctx, tx, userID, audit.Change{
		Action: ActionSettingChanged, TargetKind: "setting", TargetID: key,
		Details: map[string]json.RawMessage{"before": before, "after": value},
	})
}

func (s Service) Clear(ctx context.Context, tx pgx.Tx, userID, key string) error {
	if _, err := s.catalog.Lookup(key); err != nil {
		return err
	}
	if err := s.settings.Delete(ctx, tx, key); err != nil {
		return err
	}
	return s.log.Record(ctx, tx, userID, audit.Change{Action: ActionSettingReset, TargetKind: "setting", TargetID: key})
}

func (s Service) check(key string, value json.RawMessage) error {
	d, err := s.catalog.Lookup(key)
	if err != nil {
		return err
	}
	return d.Validate(value)
}

func (s Service) Modules(ctx context.Context, p auth.Principal) ([]ModuleState, error) {
	var off map[string]bool
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		off, err = s.modules.Disabled(ctx, tx)
		return err
	})
	if err != nil {
		return nil, err
	}
	states := []ModuleState{}
	for _, m := range s.registry.Manifests() {
		states = append(states, ModuleState{
			Name:     m.Name,
			Version:  m.Version,
			Optional: m.Optional,
			Enabled:  !m.Optional || !off[m.Name],
		})
	}
	return states, nil
}

func (s Service) SetModule(ctx context.Context, p auth.Principal, name string, enabled bool) error {
	m, ok := s.registry.Find(name)
	if !ok {
		return ErrUnknownModule
	}
	if !m.Optional && !enabled {
		return ErrRequiredModule
	}
	return s.managed(ctx, p, func(tx pgx.Tx) error {
		if err := s.modules.Set(ctx, tx, name, enabled, p.UserID); err != nil {
			return err
		}
		action := ActionModuleDisabled
		if enabled {
			action = ActionModuleEnabled
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: action, TargetKind: "module", TargetID: name})
	})
}

func (s Service) ModuleEnabled(ctx context.Context, tenantID, name string) (bool, error) {
	m, ok := s.registry.Find(name)
	if !ok {
		return false, nil
	}
	if !m.Optional {
		return true, nil
	}
	var enabled bool
	err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		var err error
		enabled, err = s.modules.Enabled(ctx, tx, name)
		return err
	})
	return enabled, err
}

func (s Service) managed(ctx context.Context, p auth.Principal, fn func(pgx.Tx) error) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), ManageSettings); err != nil {
			return err
		}
		return fn(tx)
	})
}
