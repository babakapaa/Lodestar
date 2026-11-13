package team

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct{}

func (PrivacySection) Name() string {
	return "team.teams"
}

func (PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	teams, err := tx.Query(ctx, `
		SELECT t.name, t.scope_kind, t.scope_id, m.label, m.joined_at, t.archived_at
		FROM team.members m JOIN team.teams t ON t.id = m.team_id
		WHERE m.user_id::text = $1 ORDER BY m.joined_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export teams: %w", err)
	}
	memberships, err := pgx.CollectRows(teams, func(row pgx.CollectableRow) (map[string]any, error) {
		var name, kind, id, label string
		var joined time.Time
		var archived *time.Time
		err := row.Scan(&name, &kind, &id, &label, &joined, &archived)
		return map[string]any{"team": name, "scopeKind": kind, "scopeId": id, "label": label, "joinedAt": joined, "archivedAt": archived}, err
	})
	if err != nil {
		return nil, fmt.Errorf("export teams: %w", err)
	}
	events, err := tx.Query(ctx, `
		SELECT t.name, h.event, h.at FROM team.history h JOIN team.teams t ON t.id = h.team_id
		WHERE h.user_id::text = $1 ORDER BY h.at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export team history: %w", err)
	}
	history, err := pgx.CollectRows(events, func(row pgx.CollectableRow) (map[string]any, error) {
		var name, event string
		var at time.Time
		err := row.Scan(&name, &event, &at)
		return map[string]any{"team": name, "event": event, "at": at}, err
	})
	if err != nil {
		return nil, fmt.Errorf("export team history: %w", err)
	}
	requests, err := tx.Query(ctx, `
		SELECT t.name, r.reason, r.status, r.created_at FROM team.leave_requests r JOIN team.teams t ON t.id = r.team_id
		WHERE r.user_id::text = $1 ORDER BY r.created_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export leave requests: %w", err)
	}
	leaves, err := pgx.CollectRows(requests, func(row pgx.CollectableRow) (map[string]any, error) {
		var name, reason, status string
		var at time.Time
		err := row.Scan(&name, &reason, &status, &at)
		return map[string]any{"team": name, "reason": reason, "status": status, "createdAt": at}, err
	})
	if err != nil {
		return nil, fmt.Errorf("export leave requests: %w", err)
	}
	notes, err := tx.Query(ctx, `SELECT scope_kind, scope_id, note, created_at FROM team.seekers WHERE user_id::text = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("export seeking: %w", err)
	}
	seeking, err := pgx.CollectRows(notes, func(row pgx.CollectableRow) (map[string]any, error) {
		var kind, id, note string
		var at time.Time
		err := row.Scan(&kind, &id, &note, &at)
		return map[string]any{"scopeKind": kind, "scopeId": id, "note": note, "since": at}, err
	})
	if err != nil {
		return nil, fmt.Errorf("export seeking: %w", err)
	}
	card, err := CardStore{}.Get(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	cast, err := tx.Query(ctx, `
		SELECT t.name, v.agreed, v.comment, v.at FROM team.votes v
		JOIN team.proposals p ON p.id = v.proposal_id JOIN team.teams t ON t.id = p.team_id
		WHERE v.user_id::text = $1 ORDER BY v.at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export votes: %w", err)
	}
	votes, err := pgx.CollectRows(cast, func(row pgx.CollectableRow) (map[string]any, error) {
		var name, comment string
		var agreed bool
		var at time.Time
		err := row.Scan(&name, &agreed, &comment, &at)
		return map[string]any{"team": name, "agreed": agreed, "comment": comment, "at": at}, err
	})
	if err != nil {
		return nil, fmt.Errorf("export votes: %w", err)
	}
	return map[string]any{
		"teams": memberships, "history": history, "leaveRequests": leaves, "seeking": seeking, "card": card, "votes": votes,
	}, nil
}
