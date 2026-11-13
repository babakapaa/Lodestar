package league

import (
	"cmp"
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

func (s Service) View(ctx context.Context, p auth.Principal, kind, id string) (LeagueView, error) {
	var v LeagueView
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, source, role, err := s.enter(ctx, tx, p, kind, id)
		if err != nil {
			return err
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		v, err = s.view(ctx, tx, p, role, at)
		return err
	})
	return v, err
}

func (s Service) reload(ctx context.Context, tx pgx.Tx, p auth.Principal, at synced) (LeagueView, error) {
	source, ok := s.scores.Find(at.roster.Kind)
	if !ok {
		return LeagueView{}, ErrLeagueNotFound
	}
	next, err := s.sync(ctx, tx, p.TenantID, at.roster, source)
	if err != nil {
		return LeagueView{}, err
	}
	return s.view(ctx, tx, p, RoleLead, next)
}

func (s Service) view(ctx context.Context, tx pgx.Tx, p auth.Principal, role Role, at synced) (LeagueView, error) {
	weights, err := s.weights.All(ctx, tx, at.league.ID)
	if err != nil {
		return LeagueView{}, err
	}
	v := LeagueView{
		ID: at.league.ID, Kind: at.roster.Kind, OriginID: at.roster.ID, Title: at.roster.Title, Role: role,
		Season: at.season, Divisions: at.league.Divisions, Confirm: at.league.Confirm, Formula: weights.formula(at.works),
	}
	if v.Tasks, err = s.tasks.ByLeague(ctx, tx, at.league.ID); err != nil {
		return v, err
	}
	board, err := s.board(ctx, tx, at)
	if err != nil {
		return v, err
	}
	if role == RoleLead {
		return s.lead(ctx, tx, v, at, board)
	}
	return s.student(ctx, tx, p, v, at, board)
}

func (s Service) lead(ctx context.Context, tx pgx.Tx, v LeagueView, at synced, b board) (LeagueView, error) {
	var err error
	if v.Grants, err = s.grants.Of(ctx, tx, at.league.ID, nil); err != nil {
		return v, err
	}
	v.Players = make([]PlayerRow, 0, len(b.rows))
	for _, row := range b.rows {
		v.Players = append(v.Players, PlayerRow{
			Person: b.people[row.UserID], Division: row.Division, Points: row.Points, Growth: row.Growth,
			TopGrowth: row.TopGrowth, Display: row.Display, Nickname: row.Nickname,
		})
	}
	slices.SortStableFunc(v.Players, func(a, c PlayerRow) int {
		return cmp.Or(cmp.Compare(c.Division, a.Division), cmp.Compare(c.Points, a.Points), cmp.Compare(a.Person.DisplayName, c.Person.DisplayName))
	})
	pending, err := s.transitions.Pending(ctx, tx, at.league.ID)
	if err != nil {
		return v, err
	}
	v.Pending = make([]TransitionRow, 0, len(pending))
	for _, t := range pending {
		v.Pending = append(v.Pending, TransitionRow{Transition: t, Person: b.person(t.UserID)})
	}
	open, err := s.disputes.Open(ctx, tx, at.league.ID)
	if err != nil {
		return v, err
	}
	v.Open = make([]DisputeRow, 0, len(open))
	for _, d := range open {
		row, err := s.disputeRow(ctx, tx, at.league.ID, d, b.person(d.UserID))
		if err != nil {
			return v, err
		}
		v.Open = append(v.Open, row)
	}
	return v, nil
}

func (s Service) student(ctx context.Context, tx pgx.Tx, p auth.Principal, v LeagueView, at synced, b board) (LeagueView, error) {
	me, ok := b.find(p.UserID)
	if !ok {
		return v, ErrPlayerNotFound
	}
	v.Me = &Standing{
		Points: me.Points, Division: me.Division, ToNext: at.league.Divisions.toNext(me.Division, me.Points),
		Growth: me.Growth, Display: me.Display, Nickname: me.Nickname, GPUHours: at.league.Divisions[me.Division].GPUHours,
	}
	v.Table = b.table(me.Division, p.UserID)
	var err error
	if v.Grants, err = s.grants.Of(ctx, tx, at.league.ID, &p.UserID); err != nil {
		return v, err
	}
	if v.Entries, err = s.entries.Of(ctx, tx, at.season.ID, p.UserID); err != nil {
		return v, err
	}
	if v.Transitions, err = s.transitions.Of(ctx, tx, at.league.ID, p.UserID); err != nil {
		return v, err
	}
	v.Disputes, err = s.disputes.Of(ctx, tx, at.league.ID, p.UserID)
	return v, err
}

func (s Service) disputeRow(ctx context.Context, tx pgx.Tx, leagueID string, d Dispute, person identity.Person) (DisputeRow, error) {
	row := DisputeRow{Dispute: d, Person: person}
	if d.EntryID != nil {
		e, err := s.entries.ByID(ctx, tx, leagueID, *d.EntryID)
		if err != nil {
			return row, err
		}
		row.Entry = &e
	}
	if d.TransitionID != nil {
		t, _, err := s.transitions.ByID(ctx, tx, *d.TransitionID)
		if err != nil {
			return row, err
		}
		row.Transition = &t
	}
	return row, nil
}
