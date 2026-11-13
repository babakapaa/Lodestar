package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolation = "23505"

const disputeColumns = `d.id::text, d.user_id::text, d.entry_id::text, d.transition_id::text, d.text, d.status, d.answer, d.decided_at, d.created_at`

type DisputeStore struct{}

func (DisputeStore) Insert(ctx context.Context, tx pgx.Tx, leagueID, userID string, in DisputeInput) (Dispute, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO league.disputes AS d (league_id, user_id, entry_id, transition_id, text)
		VALUES ($1, $2, $3, $4, $5) RETURNING `+disputeColumns, leagueID, userID, in.EntryID, in.TransitionID, in.Text)
	if err != nil {
		return Dispute{}, fmt.Errorf("open dispute: %w", err)
	}
	d, err := pgx.CollectExactlyOneRow(rows, scanDispute)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return Dispute{}, ErrDisputeOpen
	}
	return d, err
}

func (DisputeStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Dispute, string, error) {
	var d Dispute
	var leagueID string
	err := tx.QueryRow(ctx, `SELECT `+disputeColumns+`, d.league_id::text FROM league.disputes d WHERE d.id::text = $1 FOR UPDATE`, id).Scan(
		&d.ID, &d.UserID, &d.EntryID, &d.TransitionID, &d.Text, &d.Status, &d.Answer, &d.DecidedAt, &d.CreatedAt, &leagueID)
	if errors.Is(err, pgx.ErrNoRows) {
		return d, "", ErrDisputeNotFound
	}
	if err != nil {
		return d, "", fmt.Errorf("load dispute: %w", err)
	}
	return d, leagueID, nil
}

func (DisputeStore) Decide(ctx context.Context, tx pgx.Tx, id string, status DisputeStatus, answer, deciderID string) (Dispute, error) {
	rows, err := tx.Query(ctx, `
		UPDATE league.disputes d SET status = $2, answer = $3, decided_by = $4, decided_at = now()
		WHERE d.id::text = $1 RETURNING `+disputeColumns, id, status, answer, deciderID)
	if err != nil {
		return Dispute{}, fmt.Errorf("answer dispute: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanDispute)
}

func (DisputeStore) Open(ctx context.Context, tx pgx.Tx, leagueID string) ([]Dispute, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+disputeColumns+` FROM league.disputes d
		WHERE d.league_id::text = $1 AND d.status = 'open' ORDER BY d.created_at, d.id`, leagueID)
	if err != nil {
		return nil, fmt.Errorf("open disputes: %w", err)
	}
	return pgx.CollectRows(rows, scanDispute)
}

func (DisputeStore) Of(ctx context.Context, tx pgx.Tx, leagueID, userID string) ([]Dispute, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+disputeColumns+` FROM league.disputes d
		WHERE d.league_id::text = $1 AND d.user_id::text = $2 ORDER BY d.created_at, d.id`, leagueID, userID)
	if err != nil {
		return nil, fmt.Errorf("list disputes: %w", err)
	}
	return pgx.CollectRows(rows, scanDispute)
}

func scanDispute(row pgx.CollectableRow) (Dispute, error) {
	var d Dispute
	err := row.Scan(&d.ID, &d.UserID, &d.EntryID, &d.TransitionID, &d.Text, &d.Status, &d.Answer, &d.DecidedAt, &d.CreatedAt)
	return d, err
}
