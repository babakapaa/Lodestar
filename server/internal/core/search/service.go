package search

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const (
	defaultLimit = 20
	maxLimit     = 50
	threshold    = "0.4"
)

var kindPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)*$`)

type Modules interface {
	ModuleEnabled(ctx context.Context, tenantID, name string) (bool, error)
}

type Service struct {
	tenants database.TenantTx
	modules Modules
	sources []Source
}

func NewService(tenants database.TenantTx, modules Modules, sources ...Source) (Service, error) {
	seen := map[string]bool{}
	for _, s := range sources {
		if !kindPattern.MatchString(s.Kind()) || seen[s.Kind()] {
			return Service{}, fmt.Errorf("%w: %q", ErrBadSource, s.Kind())
		}
		seen[s.Kind()] = true
	}
	return Service{tenants: tenants, modules: modules, sources: sources}, nil
}

func (s Service) Kinds(ctx context.Context, p auth.Principal) ([]string, error) {
	on, err := s.enabled(ctx, p.TenantID, s.sources)
	if err != nil {
		return nil, err
	}
	kinds := []string{}
	for _, src := range on {
		kinds = append(kinds, src.Kind())
	}
	return kinds, nil
}

func (s Service) Search(ctx context.Context, p auth.Principal, raw string, kinds []string, limit int) ([]Hit, error) {
	q, err := ParseQuery(raw)
	if err != nil {
		return nil, err
	}
	selected, err := s.pick(kinds)
	if err != nil {
		return nil, err
	}
	if selected, err = s.enabled(ctx, p.TenantID, selected); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}
	var hits []Hit
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `SELECT set_config('pg_trgm.word_similarity_threshold', $1, true)`, threshold); err != nil {
			return err
		}
		for _, src := range selected {
			found, err := src.Search(ctx, tx, p, q, limit)
			if err != nil {
				return fmt.Errorf("search %s: %w", src.Kind(), err)
			}
			for i := range found {
				found[i].Kind = src.Kind()
			}
			hits = append(hits, found...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	slices.SortStableFunc(hits, func(a, b Hit) int {
		switch {
		case a.Score > b.Score:
			return -1
		case a.Score < b.Score:
			return 1
		}
		return 0
	})
	return hits[:min(len(hits), limit)], nil
}

func (s Service) pick(kinds []string) ([]Source, error) {
	if len(kinds) == 0 {
		return s.sources, nil
	}
	var selected []Source
	for _, kind := range kinds {
		i := slices.IndexFunc(s.sources, func(src Source) bool { return src.Kind() == kind })
		if i < 0 {
			return nil, fmt.Errorf("%w: %s", ErrUnknownKind, kind)
		}
		if !slices.Contains(selected, s.sources[i]) {
			selected = append(selected, s.sources[i])
		}
	}
	return selected, nil
}

func (s Service) enabled(ctx context.Context, tenantID string, sources []Source) ([]Source, error) {
	var on []Source
	for _, src := range sources {
		name, _, fromModule := strings.Cut(src.Kind(), ".")
		if fromModule {
			ok, err := s.modules.ModuleEnabled(ctx, tenantID, name)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
		}
		on = append(on, src)
	}
	return on, nil
}
