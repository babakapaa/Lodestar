package league

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) Awards(ctx context.Context, p auth.Principal) ([]Award, error) {
	var list []Award
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.awards.All(ctx, tx)
		return err
	})
	return list, err
}

func (s Service) CreateAward(ctx context.Context, p auth.Principal, in AwardInput) (Award, error) {
	in, err := in.normalized()
	if err != nil {
		return Award{}, err
	}
	var a Award
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), ManageAwards); err != nil {
			return err
		}
		if a, err = s.awards.Insert(ctx, tx, p.UserID, in); err != nil {
			return err
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionAward, TargetKind: "league.award", TargetID: a.ID, Details: in})
	})
	return a, err
}

func (s Service) UpdateAward(ctx context.Context, p auth.Principal, id string, in AwardInput) (Award, error) {
	in, err := in.normalized()
	if err != nil {
		return Award{}, err
	}
	var a Award
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), ManageAwards); err != nil {
			return err
		}
		if a, err = s.awards.Update(ctx, tx, id, in); err != nil {
			return err
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionAward, TargetKind: "league.award", TargetID: a.ID, Details: in})
	})
	return a, err
}

func (s Service) Grant(ctx context.Context, p auth.Principal, leagueID string, in GrantInput) (Grant, error) {
	in, err := in.normalized()
	if err != nil {
		return Grant{}, err
	}
	var g Grant
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, source, err := s.open(ctx, tx, p, leagueID, RoleLead)
		if err != nil {
			return err
		}
		if r.Archived {
			return ErrLeagueArchived
		}
		if !slices.Contains(r.Members, in.UserID) {
			return ErrPlayerNotFound
		}
		a, err := s.awards.ByID(ctx, tx, in.AwardID)
		if err != nil {
			return err
		}
		if !a.left() {
			return ErrAwardGone
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		if g, err = s.grants.Insert(ctx, tx, at.league.ID, at.season.ID, p.UserID, in); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionGrant, leagueID, map[string]any{"awardId": a.ID, "userId": in.UserID, "reason": in.Reason}); err != nil {
			return err
		}
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyAwarded, UserIDs: []string{in.UserID}, Params: map[string]any{"league": r.Title, "award": a.Title}, Link: leagueLink(r),
		})
	})
	return g, err
}

func (s Service) Revoke(ctx context.Context, p auth.Principal, id string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		leagueID, err := s.grants.League(ctx, tx, id)
		if err != nil {
			return err
		}
		if _, _, err := s.open(ctx, tx, p, leagueID, RoleLead); err != nil {
			return ErrGrantNotFound
		}
		if err := s.grants.Delete(ctx, tx, id); err != nil {
			return err
		}
		return s.record(ctx, tx, p, ActionGrant, leagueID, map[string]any{"grantId": id, "revoked": true})
	})
}
