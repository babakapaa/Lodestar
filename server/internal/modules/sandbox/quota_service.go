package sandbox

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) Quota(ctx context.Context, p auth.Principal) (Quota, error) {
	var q Quota
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, err := s.current(ctx, tx)
		if errors.Is(err, ErrNoTerm) {
			q = Quota{Lines: []QuotaLine{}, Teams: []TeamPool{}}
			return nil
		}
		if err != nil {
			return err
		}
		q, err = s.quota(ctx, tx, t, p.UserID)
		return err
	})
	return q, err
}

func (s Service) quota(ctx context.Context, tx pgx.Tx, t Term, userID string) (Quota, error) {
	lines, err := s.userLines(ctx, tx, t, userID)
	if err != nil {
		return Quota{}, err
	}
	teams, err := s.pools(ctx, tx, t, userID)
	if err != nil {
		return Quota{}, err
	}
	reserved, err := s.runs.Reserved(ctx, tx, t.ID, OwnerUser, userID)
	if err != nil {
		return Quota{}, err
	}
	q := Quota{Term: &t, Lines: lines, Hours: total(lines), Reserved: reserved, Teams: teams}
	q.Spent, q.Forecast = pace(t, lines, time.Now())
	return q, nil
}

func pace(t Term, lines []QuotaLine, now time.Time) (float64, float64) {
	var spent float64
	for _, l := range lines {
		if l.Source == SourceRun {
			spent -= l.Hours
		}
	}
	from, to := t.bounds()
	elapsed := max(now.Sub(from).Hours()/24, 1)
	left := max(to.Sub(now).Hours()/24, 0)
	return hours(spent), hours(total(lines) - spent/elapsed*left)
}

func (s Service) userLines(ctx context.Context, tx pgx.Tx, t Term, userID string) ([]QuotaLine, error) {
	base, err := s.baseHours(ctx, tx)
	if err != nil {
		return nil, err
	}
	lines := []QuotaLine{{Source: SourceBase, Hours: base}}
	from, to := t.bounds()
	grants, err := s.grants.Grants(ctx, tx, userID, from, to)
	if err != nil {
		return nil, err
	}
	for _, g := range grants {
		lines = append(lines, QuotaLine{Source: g.Source, Title: g.Title, Hours: g.Hours})
	}
	entries, err := s.ledger.Lines(ctx, tx, t.ID, OwnerUser, userID)
	if err != nil {
		return nil, err
	}
	return append(lines, entries...), nil
}

func (s Service) pools(ctx context.Context, tx pgx.Tx, t Term, userID string) ([]TeamPool, error) {
	var refs []module.TeamRef
	if s.teams != nil {
		var err error
		if refs, err = s.teams.Teams(ctx, tx, userID); err != nil {
			return nil, err
		}
	}
	pools := make([]TeamPool, 0, len(refs))
	for _, ref := range refs {
		lines, err := s.ledger.Lines(ctx, tx, t.ID, OwnerTeam, ref.ID)
		if err != nil {
			return nil, err
		}
		reserved, err := s.runs.Reserved(ctx, tx, t.ID, OwnerTeam, ref.ID)
		if err != nil {
			return nil, err
		}
		pools = append(pools, TeamPool{ID: ref.ID, Name: ref.Name, Hours: total(lines), Reserved: reserved, Lines: lines})
	}
	return pools, nil
}

func (s Service) team(ctx context.Context, tx pgx.Tx, userID, teamID string) (module.TeamRef, error) {
	if s.teams == nil {
		return module.TeamRef{}, ErrNotInTeam
	}
	refs, err := s.teams.Teams(ctx, tx, userID)
	if err != nil {
		return module.TeamRef{}, err
	}
	i := slices.IndexFunc(refs, func(r module.TeamRef) bool { return r.ID == teamID })
	if i < 0 {
		return module.TeamRef{}, ErrNotInTeam
	}
	return refs[i], nil
}

func total(lines []QuotaLine) float64 {
	var sum float64
	for _, l := range lines {
		sum += l.Hours
	}
	return hours(sum)
}
