package study

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const slotColumns = `s.id::text, s.assignment_id::text, s.starts_at, s.minutes, s.place, s.capacity,
	(SELECT count(*) FROM study.defenses d WHERE d.slot_id = s.id)::int`

type SlotStore struct{}

func (SlotStore) Insert(ctx context.Context, tx pgx.Tx, assignmentID, authorID string, in SlotsInput) ([]DefenseSlot, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO study.defense_slots AS s (assignment_id, starts_at, minutes, place, capacity, created_by)
		SELECT $1, t, $3, $4, $5, $6 FROM unnest($2::timestamptz[]) AS t
		RETURNING `+slotColumns, assignmentID, in.starts(), in.Minutes, in.Place, in.Capacity, authorID)
	if err != nil {
		return nil, fmt.Errorf("add defense slots: %w", err)
	}
	return pgx.CollectRows(rows, scanSlot)
}

func (SlotStore) ByAssignment(ctx context.Context, tx pgx.Tx, assignmentID string) ([]DefenseSlot, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+slotColumns+` FROM study.defense_slots s
		WHERE s.assignment_id::text = $1 ORDER BY s.starts_at, s.id`, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("list defense slots: %w", err)
	}
	return pgx.CollectRows(rows, scanSlot)
}

func (SlotStore) ByID(ctx context.Context, tx pgx.Tx, id string, lock bool) (DefenseSlot, error) {
	if lock {
		if _, err := tx.Exec(ctx, `SELECT 1 FROM study.defense_slots WHERE id::text = $1 FOR UPDATE`, id); err != nil {
			return DefenseSlot{}, fmt.Errorf("lock defense slot: %w", err)
		}
	}
	rows, err := tx.Query(ctx, `SELECT `+slotColumns+` FROM study.defense_slots s WHERE s.id::text = $1`, id)
	if err != nil {
		return DefenseSlot{}, fmt.Errorf("find defense slot: %w", err)
	}
	s, err := pgx.CollectExactlyOneRow(rows, scanSlot)
	if errors.Is(err, pgx.ErrNoRows) {
		return DefenseSlot{}, ErrSlotNotFound
	}
	return s, err
}

func (SlotStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM study.defense_slots WHERE id::text = $1`, id); err != nil {
		return fmt.Errorf("delete defense slot: %w", err)
	}
	return nil
}

func scanSlot(row pgx.CollectableRow) (DefenseSlot, error) {
	var s DefenseSlot
	err := row.Scan(&s.ID, &s.AssignmentID, &s.StartsAt, &s.Minutes, &s.Place, &s.Capacity, &s.Taken)
	return s, err
}
