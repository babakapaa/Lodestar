package team

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

type HistoryStore struct{}

func (HistoryStore) Add(ctx context.Context, tx pgx.Tx, teamID, userID string, event HistoryEvent, actorID string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO team.history (team_id, user_id, event, actor_id)
		VALUES ($1::uuid, nullif($2, '')::uuid, $3, nullif($4, '')::uuid)`, teamID, userID, event, actorID)
	if err != nil {
		return fmt.Errorf("team history: %w", err)
	}
	return nil
}

func (HistoryStore) ByTeam(ctx context.Context, tx pgx.Tx, teamID string) ([]HistoryEntry, error) {
	rows, err := tx.Query(ctx, `
		SELECT u.id::text, u.display_name, u.email, h.event, a.id::text, a.display_name, a.email, h.at
		FROM team.history h
		LEFT JOIN core.users u ON u.id = h.user_id
		LEFT JOIN core.users a ON a.id = h.actor_id
		WHERE h.team_id::text = $1 ORDER BY h.at, h.id`, teamID)
	if err != nil {
		return nil, fmt.Errorf("team history: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (HistoryEntry, error) {
		var e HistoryEntry
		var userID, userName, userEmail, actorID, actorName, actorEmail *string
		err := row.Scan(&userID, &userName, &userEmail, &e.Event, &actorID, &actorName, &actorEmail, &e.At)
		e.Person = person(userID, userName, userEmail)
		e.Actor = person(actorID, actorName, actorEmail)
		return e, err
	})
}

func person(id, name, email *string) *identity.Person {
	if id == nil {
		return nil
	}
	return &identity.Person{UserID: *id, DisplayName: *name, Email: *email}
}
