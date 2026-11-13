package sandbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const runColumns = `r.id::text, r.user_id::text, u.display_name, r.team_id, r.team_name, r.kind, r.title, r.command,
	r.template_id::text, t.name, t.image, r.class_id::text, c.name, r.rate::float8, r.max_hours::float8,
	r.work_kind, r.work_id, r.work_title, r.due_at, r.status, r.url, r.exit_code, r.reason, r.charged::float8,
	r.queued_at, r.started_at, r.ended_at, r.node_id::text, r.term_id::text`

const runFrom = ` FROM sandbox.runs r
	JOIN core.users u ON u.id = r.user_id
	JOIN sandbox.templates t ON t.id = r.template_id
	JOIN sandbox.gpu_classes c ON c.id = r.class_id`

const runActive = `r.status IN ('queued', 'starting', 'running', 'stopping')`

func queueOrder(days string) string {
	return `CASE WHEN r.work_kind IS NOT NULL AND r.due_at >= now() AND r.due_at <= now() + make_interval(days => ` + days + `) THEN 0 ELSE 1 END,
		CASE WHEN r.work_kind IS NOT NULL AND r.due_at >= now() AND r.due_at <= now() + make_interval(days => ` + days + `) THEN r.due_at END,
		r.weight DESC, r.queued_at, r.id`
}

type RunStore struct{}

type newRun struct {
	termID    string
	userID    string
	teamID    *string
	teamName  string
	in        RunInput
	rate      float64
	workTitle string
	dueAt     *time.Time
	weight    float64
}

func (RunStore) Insert(ctx context.Context, tx pgx.Tx, n newRun) (string, error) {
	var workKind, workID *string
	if n.in.WorkKind != "" {
		workKind, workID = &n.in.WorkKind, &n.in.WorkID
	}
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO sandbox.runs (term_id, user_id, team_id, team_name, kind, title, command, template_id, class_id,
		                          rate, max_hours, work_kind, work_id, work_title, due_at, weight)
		VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6, $7, $8::uuid, $9::uuid, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id::text`,
		n.termID, n.userID, n.teamID, n.teamName, n.in.Kind, n.in.Title, n.in.Command, n.in.TemplateID, n.in.ClassID,
		n.rate, n.in.MaxHours, workKind, workID, n.workTitle, n.dueAt, n.weight).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("queue run: %w", err)
	}
	return id, nil
}

func (RunStore) ByID(ctx context.Context, tx pgx.Tx, id string, lock bool) (Run, error) {
	sql := `SELECT ` + runColumns + runFrom + ` WHERE r.id::text = $1`
	if lock {
		sql += ` FOR UPDATE OF r`
	}
	rows, err := tx.Query(ctx, sql, id)
	if err != nil {
		return Run{}, fmt.Errorf("load run: %w", err)
	}
	r, err := pgx.CollectExactlyOneRow(rows, scanRun)
	if errors.Is(err, pgx.ErrNoRows) {
		return Run{}, ErrRunNotFound
	}
	return r, err
}

func (RunStore) Visible(ctx context.Context, tx pgx.Tx, termID, userID string, teamIDs []string) ([]Run, error) {
	rows, err := tx.Query(ctx, `SELECT `+runColumns+runFrom+`
		WHERE r.term_id::text = $1 AND (r.user_id::text = $2 OR r.team_id = ANY($3))
		ORDER BY r.queued_at DESC, r.id LIMIT 200`, termID, userID, teamIDs)
	if err != nil {
		return nil, fmt.Errorf("list runs: %w", err)
	}
	return pgx.CollectRows(rows, scanRun)
}

func (RunStore) Positions(ctx context.Context, tx pgx.Tx, days int) (map[string]int, error) {
	rows, err := tx.Query(ctx, `
		SELECT r.id::text, row_number() OVER (PARTITION BY r.class_id ORDER BY `+queueOrder("$1")+`)::int
		FROM sandbox.runs r WHERE r.status = 'queued'`, days)
	if err != nil {
		return nil, fmt.Errorf("queue positions: %w", err)
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var id string
		var pos int
		if err := rows.Scan(&id, &pos); err != nil {
			return nil, fmt.Errorf("queue positions: %w", err)
		}
		out[id] = pos
	}
	return out, rows.Err()
}

func (RunStore) Reserved(ctx context.Context, tx pgx.Tx, termID, ownerKind, ownerID string) (float64, error) {
	var sum float64
	err := tx.QueryRow(ctx, `
		SELECT coalesce(sum(r.rate * r.max_hours), 0)::float8 FROM sandbox.runs r
		WHERE r.term_id::text = $1 AND `+runActive+` AND
		      CASE WHEN $2 = 'team' THEN r.team_id = $3 ELSE r.team_id IS NULL AND r.user_id::text = $3 END`,
		termID, ownerKind, ownerID).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("reserved hours: %w", err)
	}
	return hours(sum), nil
}

func (RunStore) QueuedClasses(ctx context.Context, tx pgx.Tx) ([]string, error) {
	rows, err := tx.Query(ctx, `SELECT DISTINCT class_id::text FROM sandbox.runs WHERE status = 'queued'`)
	if err != nil {
		return nil, fmt.Errorf("queued classes: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (RunStore) Queued(ctx context.Context, tx pgx.Tx, classID string, days, limit int) ([]Run, error) {
	rows, err := tx.Query(ctx, `SELECT `+runColumns+runFrom+`
		WHERE r.class_id::text = $2 AND r.status = 'queued' ORDER BY `+queueOrder("$1")+` LIMIT $3 FOR UPDATE OF r`,
		days, classID, limit)
	if err != nil {
		return nil, fmt.Errorf("queued runs: %w", err)
	}
	return pgx.CollectRows(rows, scanRun)
}

func (RunStore) OnNode(ctx context.Context, tx pgx.Tx, nodeID string) ([]Run, error) {
	rows, err := tx.Query(ctx, `SELECT `+runColumns+runFrom+`
		WHERE r.node_id::text = $1 AND r.status IN ('starting', 'running', 'stopping') ORDER BY r.queued_at FOR UPDATE OF r`, nodeID)
	if err != nil {
		return nil, fmt.Errorf("runs on node: %w", err)
	}
	return pgx.CollectRows(rows, scanRun)
}

func (RunStore) Overdue(ctx context.Context, tx pgx.Tx) ([]Run, error) {
	rows, err := tx.Query(ctx, `SELECT `+runColumns+runFrom+`
		WHERE r.status = 'running' AND r.started_at + make_interval(secs => r.max_hours * 3600) <= now() FOR UPDATE OF r`)
	if err != nil {
		return nil, fmt.Errorf("overdue runs: %w", err)
	}
	return pgx.CollectRows(rows, scanRun)
}

func (RunStore) Assign(ctx context.Context, tx pgx.Tx, id, nodeID string) error {
	if _, err := tx.Exec(ctx, `UPDATE sandbox.runs SET status = 'starting', node_id = $2::uuid WHERE id::text = $1`, id, nodeID); err != nil {
		return fmt.Errorf("assign run: %w", err)
	}
	return nil
}

func (RunStore) Started(ctx context.Context, tx pgx.Tx, id, url string) error {
	if _, err := tx.Exec(ctx, `
		UPDATE sandbox.runs SET status = CASE WHEN status = 'starting' THEN 'running' ELSE status END,
		       started_at = coalesce(started_at, now()), url = $2
		WHERE id::text = $1`, id, url); err != nil {
		return fmt.Errorf("start run: %w", err)
	}
	return nil
}

func (RunStore) Stopping(ctx context.Context, tx pgx.Tx, id, reason string) error {
	if _, err := tx.Exec(ctx, `UPDATE sandbox.runs SET status = 'stopping', reason = $2 WHERE id::text = $1`, id, reason); err != nil {
		return fmt.Errorf("stop run: %w", err)
	}
	return nil
}

func (RunStore) End(ctx context.Context, tx pgx.Tx, id, status string, exitCode *int, reason string, charged float64, at time.Time) error {
	if _, err := tx.Exec(ctx, `
		UPDATE sandbox.runs SET status = $2, exit_code = $3, reason = CASE WHEN $4 = '' THEN reason ELSE $4 END,
		       charged = $5, ended_at = $6, url = ''
		WHERE id::text = $1`, id, status, exitCode, reason, charged, at); err != nil {
		return fmt.Errorf("end run: %w", err)
	}
	return nil
}

func scanRun(row pgx.CollectableRow) (Run, error) {
	var r Run
	err := row.Scan(&r.ID, &r.UserID, &r.DisplayName, &r.TeamID, &r.TeamName, &r.Kind, &r.Title, &r.Command,
		&r.TemplateID, &r.TemplateName, &r.Image, &r.ClassID, &r.ClassName, &r.Rate, &r.MaxHours,
		&r.WorkKind, &r.WorkID, &r.WorkTitle, &r.DueAt, &r.Status, &r.URL, &r.ExitCode, &r.Reason, &r.Charged,
		&r.QueuedAt, &r.StartedAt, &r.EndedAt, &r.NodeID, &r.TermID)
	return r, err
}
