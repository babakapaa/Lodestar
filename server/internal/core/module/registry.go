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

var ErrForeignRoster = errors.New("module roster kind must start with the module name")

var ErrForeignTeamScope = errors.New("module team scope kind must start with the module name")

var ErrForeignScoreSource = errors.New("module score source kind must start with the module name")

var ErrForeignLinkSource = errors.New("module link source kind must start with the module name")

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9]{1,29}$`)

var reservedNames = []string{
	"auth", "events", "files", "live", "lti", "me", "meta", "modules", "notifications", "open", "overrides", "permissions",
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
		if err := checkRosterKinds(m); err != nil {
			return nil, err
		}
		if err := checkTeamScopeKinds(m); err != nil {
			return nil, err
		}
		if err := checkScoreKinds(m); err != nil {
			return nil, err
		}
		if err := checkLinkKinds(m); err != nil {
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
		if o, ok := m.(OpenMounting); ok {
			open := http.NewServeMux()
			o.MountOpen(open, k)
			mux.Handle("/v1/open/"+name+"/", open)
		}
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

func (r *Registry) LinkSources(k Kernel) LinkSources {
	var sources LinkSources
	for _, m := range r.modules {
		if s, ok := m.(LinkProviding); ok {
			sources = append(sources, s.LinkSources(k)...)
		}
	}
	return sources
}

func (r *Registry) Rosters(k Kernel) Rosters {
	var rosters Rosters
	for _, m := range r.modules {
		if s, ok := m.(Rostering); ok {
			rosters = append(rosters, s.RosterSources(k)...)
		}
	}
	return rosters
}

func (r *Registry) TeamScopes(k Kernel) TeamScopes {
	var scopes TeamScopes
	for _, m := range r.modules {
		if s, ok := m.(Teaming); ok {
			scopes = append(scopes, s.TeamScopeSources(k)...)
		}
	}
	return scopes
}

func (r *Registry) TeamDirectory(k Kernel) TeamDirectory {
	for _, m := range r.modules {
		if d, ok := m.(TeamDirecting); ok {
			return d.TeamDirectory(k)
		}
	}
	return nil
}

func (r *Registry) LessonSource(k Kernel) LessonSource {
	for _, m := range r.modules {
		if l, ok := m.(LessonProviding); ok {
			return l.LessonSource(k)
		}
	}
	return nil
}

func (r *Registry) StatsSources(k Kernel) []StatsSource {
	var sources []StatsSource
	for _, m := range r.modules {
		if s, ok := m.(StatsProviding); ok {
			sources = append(sources, s.StatsSources(k)...)
		}
	}
	return sources
}

func (r *Registry) ReadinessSources(k Kernel) []ReadinessSource {
	var sources []ReadinessSource
	for _, m := range r.modules {
		if s, ok := m.(ReadinessProviding); ok {
			sources = append(sources, s.ReadinessSources(k)...)
		}
	}
	return sources
}

func (r *Registry) ScoreSources(k Kernel) ScoreSources {
	var sources ScoreSources
	for _, m := range r.modules {
		if s, ok := m.(Scoring); ok {
			sources = append(sources, s.ScoreSources(k)...)
		}
	}
	return sources
}

func (r *Registry) GPUSources(k Kernel) GPUSources {
	var sources GPUSources
	for _, m := range r.modules {
		if g, ok := m.(GPUGranting); ok {
			sources = append(sources, g.GPUSources(k)...)
		}
	}
	return sources
}
