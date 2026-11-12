package module

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

var ErrDuplicateModule = errors.New("module registered twice")

var ErrUnnamedModule = errors.New("module has no name")

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

func (r *Registry) Mount(mux *http.ServeMux) {
	for _, m := range r.modules {
		m.Mount(mux)
	}
}
