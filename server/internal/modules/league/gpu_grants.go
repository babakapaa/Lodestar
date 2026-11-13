package league

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const (
	GPUSourceDivision = "league.division"
	GPUSourceAward    = "league.award"
)

type GPUGrants struct{}

func (GPUGrants) GPUGrants(ctx context.Context, tx pgx.Tx, userID string, from, to time.Time) ([]module.GPUGrant, error) {
	var out []module.GPUGrant
	var division string
	var hours int
	err := tx.QueryRow(ctx, `
		SELECT d->>'name', (d->>'gpuHours')::int
		FROM league.players p
		JOIN league.leagues l ON l.id = p.league_id
		JOIN league.seasons s ON s.league_id = l.id AND s.ended_at IS NULL
		CROSS JOIN LATERAL (SELECT l.divisions->p.division AS d) x
		WHERE p.user_id::text = $1 AND (d->>'gpuHours')::int > 0
		ORDER BY (d->>'gpuHours')::int DESC LIMIT 1`, userID).Scan(&division, &hours)
	switch {
	case err == nil:
		out = append(out, module.GPUGrant{Source: GPUSourceDivision, Title: division, Hours: float64(hours), Priority: true})
	case !errors.Is(err, pgx.ErrNoRows):
		return nil, fmt.Errorf("division gpu hours: %w", err)
	}
	rows, err := tx.Query(ctx, `
		SELECT a.title, a.gpu_hours FROM league.grants g JOIN league.awards a ON a.id = g.award_id
		WHERE g.user_id::text = $1 AND a.gpu_hours > 0 AND g.granted_at >= $2 AND g.granted_at < $3
		ORDER BY g.granted_at`, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("award gpu hours: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		g := module.GPUGrant{Source: GPUSourceAward}
		var awardHours int
		if err := rows.Scan(&g.Title, &awardHours); err != nil {
			return nil, fmt.Errorf("award gpu hours: %w", err)
		}
		g.Hours = float64(awardHours)
		out = append(out, g)
	}
	return out, rows.Err()
}
