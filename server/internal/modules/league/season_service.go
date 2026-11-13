package league

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) CloseSeason(ctx context.Context, p auth.Principal, leagueID string) (LeagueView, error) {
	var v LeagueView
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, source, err := s.open(ctx, tx, p, leagueID, RoleLead)
		if err != nil {
			return err
		}
		if r.Archived {
			return ErrLeagueArchived
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		totals, err := s.entries.Totals(ctx, tx, at.season.ID)
		if err != nil {
			return err
		}
		players, err := s.players.All(ctx, tx, at.league.ID)
		if err != nil {
			return err
		}
		divisions := at.league.Divisions
		moved, queued := 0, false
		for _, userID := range r.Members {
			from := divisions.clamp(players[userID].Division)
			to := divisions.at(totals[userID])
			if to >= from {
				continue
			}
			t := Transition{SeasonID: at.season.ID, UserID: userID, From: from, To: to, Points: totals[userID]}
			created, err := s.shift(ctx, tx, p.TenantID, at, t)
			if err != nil {
				return err
			}
			moved++
			queued = queued || created
		}
		if queued {
			if err := s.tellQueued(ctx, tx, p.TenantID, r); err != nil {
				return err
			}
		}
		if err := s.seasons.End(ctx, tx, at.season.ID); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionSeason, at.league.ID, map[string]any{"season": at.season.Number, "down": moved}); err != nil {
			return err
		}
		v, err = s.reload(ctx, tx, p, at)
		return err
	})
	return v, err
}
