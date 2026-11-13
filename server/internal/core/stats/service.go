package stats

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const (
	Threshold   = 5
	defaultDays = 30
	longestDays = 400
)

type Service struct {
	tenants   database.TenantTx
	access    access.Service
	sources   []module.StatsSource
	units     structure.UnitStore
	directory structure.Directory
	now       func() time.Time
}

func NewService(tenants database.TenantTx, acc access.Service, sources ...module.StatsSource) Service {
	return Service{tenants: tenants, access: acc, sources: sources, now: time.Now}
}

func (s Service) Units(ctx context.Context, p auth.Principal, parentID *string, days int) (Report, error) {
	if days == 0 {
		days = defaultDays
	}
	if days < 1 || days > longestDays {
		return Report{}, ErrInvalidPeriod
	}
	out := Report{ParentID: parentID, Days: days, Threshold: Threshold, Units: []UnitStats{}}
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		scope := access.TenantScope(p.TenantID)
		if parentID != nil {
			if _, err := s.units.ByID(ctx, tx, *parentID); err != nil {
				return err
			}
			scope = structure.UnitScope(*parentID)
		}
		full, err := s.allowed(ctx, tx, p, scope, Full)
		if err != nil {
			return err
		}
		if !full {
			if err := s.access.Require(ctx, tx, p, scope, View); err != nil {
				return err
			}
		}
		out.Full = full
		children, err := s.units.Children(ctx, tx, parentID, false)
		if err != nil {
			return err
		}
		ids := make([]string, len(children))
		for i, u := range children {
			ids[i] = u.ID
		}
		people, err := s.directory.UsersUnderEach(ctx, tx, ids)
		if err != nil {
			return err
		}
		to := s.now()
		from := to.AddDate(0, 0, -days)
		for _, u := range children {
			members := people[u.ID]
			row := UnitStats{UnitID: u.ID, Kind: u.Kind, Name: u.Name, Metrics: []module.Metric{}}
			if len(members) < Threshold && !full {
				row.Hidden = true
				out.Units = append(out.Units, row)
				continue
			}
			n := len(members)
			row.People = &n
			for _, source := range s.sources {
				metrics, err := source.Stats(ctx, tx, members, from, to)
				if err != nil {
					return err
				}
				row.Metrics = append(row.Metrics, metrics...)
			}
			out.Units = append(out.Units, row)
		}
		return nil
	})
	return out, err
}

func (s Service) allowed(ctx context.Context, tx pgx.Tx, p auth.Principal, scope access.Scope, perm access.Permission) (bool, error) {
	err := s.access.Require(ctx, tx, p, scope, perm)
	if errors.Is(err, access.ErrForbidden) {
		return false, nil
	}
	return err == nil, err
}
