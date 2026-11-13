package league

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type synced struct {
	league League
	season Season
	works  []module.ScoredWork
	roster module.Roster
}

func (s Service) sync(ctx context.Context, tx pgx.Tx, tenantID string, r module.Roster, source module.ScoreSource) (synced, error) {
	l, err := s.leagues.Ensure(ctx, tx, r.Kind, r.ID)
	if err != nil {
		return synced{}, err
	}
	season, err := s.seasons.Current(ctx, tx, l.ID)
	if err != nil {
		return synced{}, err
	}
	works, err := source.ScoredWorks(ctx, tx, r.ID)
	if err != nil {
		return synced{}, err
	}
	out := synced{league: l, season: season, works: works, roster: r}
	if r.Archived {
		return out, nil
	}
	if err := s.players.Ensure(ctx, tx, l.ID, r.Members); err != nil {
		return out, err
	}
	weights, err := s.weights.All(ctx, tx, l.ID)
	if err != nil {
		return out, err
	}
	if err := s.score(ctx, tx, out, weights); err != nil {
		return out, err
	}
	return out, s.promote(ctx, tx, tenantID, out)
}

func (s Service) score(ctx context.Context, tx pgx.Tx, at synced, weights Weights) error {
	existing, err := s.entries.Works(ctx, tx, at.league.ID, EntryWork)
	if err != nil {
		return err
	}
	members := make(map[string]bool, len(at.roster.Members))
	for _, id := range at.roster.Members {
		members[id] = true
	}
	kept := map[string]bool{}
	for _, work := range at.works {
		weight := weights.of(work)
		for _, result := range work.Results {
			if !members[result.UserID] {
				continue
			}
			key := workKey(result.UserID, work.Kind, work.ID)
			old, found := existing[key]
			if found && old.SeasonID != at.season.ID {
				kept[key] = true
				continue
			}
			points := int(math.Round(result.Share() * float64(weight)))
			if points == 0 {
				continue
			}
			kept[key] = true
			next := Entry{
				ID: old.ID, UserID: result.UserID, Kind: EntryWork, WorkKind: &work.Kind, WorkID: &work.ID, Title: work.Title,
				Score: &result.Score, Max: &result.Max, Weight: &weight, Points: points,
			}
			switch {
			case !found:
				err = s.entries.InsertWork(ctx, tx, at.league.ID, at.season.ID, next)
			case changed(old, next):
				err = s.entries.UpdateWork(ctx, tx, next)
			}
			if err != nil {
				return err
			}
		}
	}
	for key, old := range existing {
		if !kept[key] && old.SeasonID == at.season.ID {
			if err := s.entries.Delete(ctx, tx, old.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s Service) promote(ctx context.Context, tx pgx.Tx, tenantID string, at synced) error {
	totals, err := s.entries.Totals(ctx, tx, at.season.ID)
	if err != nil {
		return err
	}
	players, err := s.players.All(ctx, tx, at.league.ID)
	if err != nil {
		return err
	}
	divisions := at.league.Divisions
	queued := false
	for _, userID := range at.roster.Members {
		from := divisions.clamp(players[userID].Division)
		to := divisions.at(totals[userID])
		if to <= from {
			continue
		}
		highest, err := s.transitions.Highest(ctx, tx, at.season.ID, userID)
		if err != nil {
			return err
		}
		if highest >= to {
			continue
		}
		t := Transition{SeasonID: at.season.ID, UserID: userID, From: from, To: to, Points: totals[userID]}
		created, err := s.shift(ctx, tx, tenantID, at, t)
		if err != nil {
			return err
		}
		queued = queued || created
	}
	if !queued {
		return nil
	}
	return s.tellQueued(ctx, tx, tenantID, at.roster)
}

func (s Service) shift(ctx context.Context, tx pgx.Tx, tenantID string, at synced, t Transition) (bool, error) {
	if at.league.Confirm {
		_, created, err := s.transitions.Queue(ctx, tx, at.league.ID, t)
		return created, err
	}
	applied, err := s.transitions.Apply(ctx, tx, at.league.ID, t)
	if err != nil {
		return false, err
	}
	if err := s.players.SetDivision(ctx, tx, at.league.ID, t.UserID, t.To); err != nil {
		return false, err
	}
	return false, s.tellMoved(ctx, tx, tenantID, at, applied)
}

func changed(old, next Entry) bool {
	return old.Title != next.Title || old.Reason != next.Reason || old.Points != next.Points || *old.Score != *next.Score || *old.Max != *next.Max || *old.Weight != *next.Weight
}
