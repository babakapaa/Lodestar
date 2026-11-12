package module

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

var ErrDuplicateModule = errors.New("module registered twice")

var ErrUnnamedModule = errors.New("module has no name")

var ErrReservedName = errors.New("module name is reserved or malformed")

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9]{1,29}$`)

var reservedNames = []string{
	"auth", "events", "live", "me", "meta", "modules", "notifications", "overrides", "permissions",
	"preferences", "role-assignments", "roles", "sessions", "settings", "units", "users",
}

type Registry struct {
	modules []Module
}

func NewRegistry(modules ...Module) (*Registry, error) {
	seen := make(map[string]struct{}, len(modules))
	for _, m := range modules {
		name := m.Manifest().Name
		if name == "" {
			return nil, ErrUnnamedModule
		}
		if !namePattern.MatchString(name) || slices.Contains(reservedNames, name) {
			return nil, fmt.Errorf("%w: %s", ErrReservedName, name)
		}
		if _, ok := seen[name]; ok {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateModule, name)
		}
		seen[name] = struct{}{}
	}
	return &Registry{modules: modules}, nil
}

func (r *Registry) Manifests() []Manifest {
	manifests := make([]Manifest, 0, len(r.modules))
	for _, m := range r.modules {
		manifests = append(manifests, m.Manifest())
	}
	return manifests
}

func (r *Registry) Migrations() []database.Migrations {
	var sets []database.Migrations
	for _, m := range r.modules {
		if migrating, ok := m.(Migrating); ok {
			sets = append(sets, migrating.Migrations())
		}
	}
	return sets
}

func (r *Registry) Permissions() []string {
	var perms []string
	for _, m := range r.modules {
		if p, ok := m.(Permissioned); ok {
			perms = append(perms, p.Permissions()...)
		}
	}
	return perms
}

func (r *Registry) Preferences() []schema.Definition {
	var defs []schema.Definition
	for _, m := range r.modules {
		if p, ok := m.(Preferring); ok {
			defs = append(defs, p.Preferences()...)
		}
	}
	return defs
}

func (r *Registry) Settings() []schema.Definition {
	var defs []schema.Definition
	for _, m := range r.modules {
		if c, ok := m.(Configurable); ok {
			defs = append(defs, c.Settings()...)
		}
	}
	return defs
}

func (r *Registry) NotificationTypes() []NotificationType {
	var types []NotificationType
	for _, m := range r.modules {
		if n, ok := m.(Notifying); ok {
			types = append(types, n.NotificationTypes()...)
		}
	}
	return types
}

func (r *Registry) Find(name string) (Manifest, bool) {
	for _, m := range r.modules {
		if m.Manifest().Name == name {
			return m.Manifest(), true
		}
	}
	return Manifest{}, false
}

func (r *Registry) Mount(mux *http.ServeMux, gate Gate) {
	for _, m := range r.modules {
		name := m.Manifest().Name
		sub := http.NewServeMux()
		m.Mount(sub)
		mux.Handle("/v1/"+name+"/", gate.Wrap(name, sub))
	}
}
