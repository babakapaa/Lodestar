package league

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) tellQueued(ctx context.Context, tx pgx.Tx, tenantID string, r module.Roster) error {
	if len(r.Leads) == 0 {
		return nil
	}
	return s.notifier.Notify(ctx, tx, tenantID, module.Notification{
		Type: NotifyQueued, UserIDs: r.Leads, Params: map[string]any{"league": r.Title}, Link: leagueLink(r),
	})
}

func (s Service) tellMoved(ctx context.Context, tx pgx.Tx, tenantID string, at synced, t Transition) error {
	kind := NotifyRelegated
	if t.up() {
		kind = NotifyPromoted
	}
	return s.notifier.Notify(ctx, tx, tenantID, module.Notification{
		Type: kind, UserIDs: []string{t.UserID},
		Params: map[string]any{"league": at.roster.Title, "division": at.league.Divisions.label(t.To)},
		Link:   leagueLink(at.roster),
	})
}

func (s Service) tellDisputed(ctx context.Context, tx pgx.Tx, tenantID string, r module.Roster, author string) error {
	if len(r.Leads) == 0 {
		return nil
	}
	return s.notifier.Notify(ctx, tx, tenantID, module.Notification{
		Type: NotifyDisputed, UserIDs: r.Leads, Params: map[string]any{"league": r.Title, "student": author}, Link: leagueLink(r),
	})
}

func (s Service) tellAnswered(ctx context.Context, tx pgx.Tx, tenantID string, r module.Roster, d Dispute) error {
	return s.notifier.Notify(ctx, tx, tenantID, module.Notification{
		Type: NotifyAnswered, UserIDs: []string{d.UserID}, Params: map[string]any{"league": r.Title}, Link: leagueLink(r),
	})
}

func leagueLink(r module.Roster) string {
	return "/league/" + r.Kind + "/" + r.ID
}
