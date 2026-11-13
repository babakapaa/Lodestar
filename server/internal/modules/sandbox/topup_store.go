package sandbox

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const topupColumns = `t.id::text, t.term_id::text, t.user_id::text, u.display_name, u.email, t.work_kind, t.work_id, t.title,
	t.hours::float8, t.reason, t.status, t.granted::float8, t.created_at, t.decided_at`

const topupFrom = ` FROM sandbox.topups t JOIN core.users u ON u.id = t.user_id`

type TopupStore struct{}

func (TopupStore) Insert(ctx context.Context, tx pgx.Tx, termID, userID, title string, in TopupInput) (Topup, error) {
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO sandbox.topups (term_id, user_id, work_kind, work_id, title, hours, reason)
		VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6, $7) RETURNING id::text`,
		termID, userID, in.WorkKind, in.WorkID, title, in.Hours, in.Reason).Scan(&id)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return Topup{}, ErrTopupPending
	}
	if err != nil {
		return Topup{}, fmt.Errorf("request top-up: %w", err)
	}
	return TopupStore{}.one(ctx, tx, `SELECT `+topupColumns+topupFrom+` WHERE t.id::text = $1`, id)
}

func (TopupStore) Pending(ctx context.Context, tx pgx.Tx, id string) (Topup, error) {
	return TopupStore{}.one(ctx, tx, `SELECT `+topupColumns+topupFrom+` WHERE t.id::text = $1 AND t.status = 'pending' FOR UPDATE OF t`, id)
}

func (TopupStore) Decide(ctx context.Context, tx pgx.Tx, id, status string, granted *float64, actorID string) (Topup, error) {
	if _, err := tx.Exec(ctx, `
		UPDATE sandbox.topups SET status = $2, granted = $3, decided_by = $4::uuid, decided_at = now()
		WHERE id::text = $1`, id, status, granted, actorID); err != nil {
		return Topup{}, fmt.Errorf("decide top-up: %w", err)
	}
	return TopupStore{}.one(ctx, tx, `SELECT `+topupColumns+topupFrom+` WHERE t.id::text = $1`, id)
}

func (TopupStore) InTerm(ctx context.Context, tx pgx.Tx, termID string) ([]Topup, error) {
	rows, err := tx.Query(ctx, `SELECT `+topupColumns+topupFrom+`
		WHERE t.term_id::text = $1 ORDER BY t.status <> 'pending', t.created_at DESC`, termID)
	if err != nil {
		return nil, fmt.Errorf("list top-ups: %w", err)
	}
	return pgx.CollectRows(rows, scanTopup)
}

func (TopupStore) one(ctx context.Context, tx pgx.Tx, sql string, id string) (Topup, error) {
	rows, err := tx.Query(ctx, sql, id)
	if err != nil {
		return Topup{}, fmt.Errorf("load top-up: %w", err)
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTopup)
	if errors.Is(err, pgx.ErrNoRows) {
		return Topup{}, ErrTopupNotFound
	}
	return t, err
}

func scanTopup(row pgx.CollectableRow) (Topup, error) {
	var t Topup
	err := row.Scan(&t.ID, &t.TermID, &t.UserID, &t.DisplayName, &t.Email, &t.WorkKind, &t.WorkID, &t.Title,
		&t.Hours, &t.Reason, &t.Status, &t.Granted, &t.CreatedAt, &t.DecidedAt)
	return t, err
}
