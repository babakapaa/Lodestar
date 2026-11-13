package league

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct{}

func (PrivacySection) Name() string {
	return "league.players"
}

func (PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	rows, err := tx.Query(ctx, `
		SELECT l.origin_kind, l.origin_id, p.division, p.display, p.nickname, p.joined_at
		FROM league.players p JOIN league.leagues l ON l.id = p.league_id
		WHERE p.user_id::text = $1 ORDER BY p.joined_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export leagues: %w", err)
	}
	leagues, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (map[string]any, error) {
		var kind, id, display, nickname string
		var division int
		var joined time.Time
		err := row.Scan(&kind, &id, &division, &display, &nickname, &joined)
		return map[string]any{"kind": kind, "originId": id, "division": division, "display": display, "nickname": nickname, "joinedAt": joined}, err
	})
	if err != nil {
		return nil, fmt.Errorf("export leagues: %w", err)
	}
	rows, err = tx.Query(ctx, `
		SELECT l.origin_id, s.number, e.kind, e.title, e.score, e.max, e.weight, e.points, e.reason, e.created_at
		FROM league.entries e JOIN league.leagues l ON l.id = e.league_id JOIN league.seasons s ON s.id = e.season_id
		WHERE e.user_id::text = $1 ORDER BY e.created_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export league points: %w", err)
	}
	points, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (map[string]any, error) {
		var origin, kind, title, reason string
		var season, value int
		var score, maxScore *float64
		var weight *int
		var at time.Time
		err := row.Scan(&origin, &season, &kind, &title, &score, &maxScore, &weight, &value, &reason, &at)
		return map[string]any{
			"originId": origin, "season": season, "kind": kind, "title": title, "score": score, "max": maxScore,
			"weight": weight, "points": value, "reason": reason, "at": at,
		}, err
	})
	if err != nil {
		return nil, fmt.Errorf("export league points: %w", err)
	}
	rows, err = tx.Query(ctx, `
		SELECT l.origin_id, t.from_division, t.to_division, t.points, t.status, t.created_at
		FROM league.transitions t JOIN league.leagues l ON l.id = t.league_id
		WHERE t.user_id::text = $1 ORDER BY t.created_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export league moves: %w", err)
	}
	moves, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (map[string]any, error) {
		var origin, status string
		var from, to, value int
		var at time.Time
		err := row.Scan(&origin, &from, &to, &value, &status, &at)
		return map[string]any{"originId": origin, "from": from, "to": to, "points": value, "status": status, "at": at}, err
	})
	if err != nil {
		return nil, fmt.Errorf("export league moves: %w", err)
	}
	rows, err = tx.Query(ctx, `
		SELECT l.origin_id, d.text, d.status, d.answer, d.created_at
		FROM league.disputes d JOIN league.leagues l ON l.id = d.league_id
		WHERE d.user_id::text = $1 ORDER BY d.created_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export league disputes: %w", err)
	}
	disputes, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (map[string]any, error) {
		var origin, text, status, answer string
		var at time.Time
		err := row.Scan(&origin, &text, &status, &answer, &at)
		return map[string]any{"originId": origin, "text": text, "status": status, "answer": answer, "at": at}, err
	})
	if err != nil {
		return nil, fmt.Errorf("export league disputes: %w", err)
	}
	return map[string]any{"leagues": leagues, "points": points, "moves": moves, "disputes": disputes}, nil
}
