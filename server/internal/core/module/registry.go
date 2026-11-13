package module

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"

	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

var ErrDuplicateModule = errors.New("module registered twice")

var ErrUnnamedModule = errors.New("module has no name")

var ErrReservedName = errors.New("module name is reserved or malformed")

var ErrForeignSearchKind = errors.New("module search kind must start with the module name")

var ErrForeignSection = errors.New("module privacy section must start with the module name")

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9]{1,29}$`)

var reservedNames = []string{
	"auth", "events", "files", "live", "me", "meta", "modules", "notifications", "overrides", "permissions",
	"preferences", "role-assignments", "roles", "search", "sessions", "settings", "units", "users",
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
		if err := checkSearchKinds(m); err != nil {
			return nil, err
		}
		if err := checkPrivacySections(m); err != nil {
			return nil, err
		}
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

func (r *Registry) SearchSources() []search.Source {
	var sources []search.Source
	for _, m := range r.modules {
		if s, ok := m.(Searchable); ok {
			sources = append(sources, s.SearchSources()...)
		}
	}
	return sources
}

func (r *Registry) PrivacySections(k Kernel) []privacy.Section {
	var sections []privacy.Section
	for _, m := range r.modules {
		if p, ok := m.(Private); ok {
			sections = append(sections, p.PrivacySections(k)...)
		}
	}
	return sections
}

func (r *Registry) FileLinks(k Kernel) []files.Link {
	var links []files.Link
	for _, m := range r.modules {
		if l, ok := m.(Linking); ok {
			links = append(links, l.FileLinks(k)...)
		}
	}
	return links
}

func (r *Registry) Workers(k Kernel) []Worker {
	var workers []Worker
	for _, m := range r.modules {
		if w, ok := m.(Working); ok {
			workers = append(workers, w.Workers(k)...)
		}
	}
	return workers
}

func (r *Registry) Find(name string) (Manifest, bool) {
	for _, m := range r.modules {
		if m.Manifest().Name == name {
			return m.Manifest(), true
		}
	}
	return Manifest{}, false
}

func (r *Registry) Mount(mux *http.ServeMux, gate Gate, k Kernel) {
	for _, m := range r.modules {
		name := m.Manifest().Name
		sub := http.NewServeMux()
		m.Mount(sub, k)
		gated := gate.Wrap(name, sub)
		mux.Handle("/v1/"+name, gated)
		mux.Handle("/v1/"+name+"/", gated)
	}
}

func (r *Registry) Briefs(k Kernel) Briefs {
	var briefs Briefs
	for _, m := range r.modules {
		if b, ok := m.(Briefing); ok {
			briefs = append(briefs, b.Briefers(k)...)
		}
	}
	return briefs
}
