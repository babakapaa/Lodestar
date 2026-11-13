package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const reportColumns = `id::text, chat_id::text, reporter_id::text, data, created_at, resolved_at`

type ReportStore struct {
	sealer Sealer
}

func (s ReportStore) Insert(ctx context.Context, tx pgx.Tx, chatID, reporterID string, body reportBody) (Report, error) {
	data, err := s.sealer.seal(body)
	if err != nil {
		return Report{}, err
	}
	rows, err := tx.Query(ctx, `
		INSERT INTO chat.reports (chat_id, reporter_id, data) VALUES ($1, $2, $3)
		RETURNING `+reportColumns, chatID, reporterID, data)
	if err != nil {
		return Report{}, fmt.Errorf("insert report: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, s.scan)
}

func (s ReportStore) Open(ctx context.Context, tx pgx.Tx) ([]Report, error) {
	rows, err := tx.Query(ctx, `SELECT `+reportColumns+` FROM chat.reports WHERE resolved_at IS NULL ORDER BY created_at, id`)
	if err != nil {
		return nil, fmt.Errorf("open reports: %w", err)
	}
	return pgx.CollectRows(rows, s.scan)
}

func (s ReportStore) ByReporter(ctx context.Context, tx pgx.Tx, reporterID string) ([]Report, error) {
	rows, err := tx.Query(ctx, `SELECT `+reportColumns+` FROM chat.reports WHERE reporter_id::text = $1 ORDER BY created_at, id`, reporterID)
	if err != nil {
		return nil, fmt.Errorf("reports by reporter: %w", err)
	}
	return pgx.CollectRows(rows, s.scan)
}

func (s ReportStore) Resolve(ctx context.Context, tx pgx.Tx, id, userID string) (Report, error) {
	rows, err := tx.Query(ctx, `
		UPDATE chat.reports SET resolved_by = $2, resolved_at = coalesce(resolved_at, now())
		WHERE id::text = $1 RETURNING `+reportColumns, id, userID)
	if err != nil {
		return Report{}, fmt.Errorf("resolve report: %w", err)
	}
	r, err := pgx.CollectExactlyOneRow(rows, s.scan)
	if errors.Is(err, pgx.ErrNoRows) {
		return Report{}, ErrReportNotFound
	}
	return r, err
}

func (s ReportStore) scan(row pgx.CollectableRow) (Report, error) {
	var r Report
	var data []byte
	if err := row.Scan(&r.ID, &r.ChatID, &r.ReporterID, &data, &r.CreatedAt, &r.ResolvedAt); err != nil {
		return Report{}, err
	}
	var body reportBody
	if err := s.sealer.open(data, &body); err != nil {
		return Report{}, err
	}
	r.ChatTitle, r.Reason, r.Messages = body.ChatTitle, body.Reason, body.Messages
	return r, nil
}
