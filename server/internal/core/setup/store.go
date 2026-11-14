package setup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type Store struct{}

func (Store) State(ctx context.Context, tx pgx.Tx) (State, bool, error) {
	var s State
	err := tx.QueryRow(ctx, `
		SELECT protocol_version, signed_by::text, signed_at, launched_at, pilot_courses::text[] FROM core.setup_state`).
		Scan(&s.ProtocolVersion, &s.SignedBy, &s.SignedAt, &s.LaunchedAt, &s.Pilot)
	if errors.Is(err, pgx.ErrNoRows) {
		return State{}, false, nil
	}
	if err != nil {
		return State{}, false, fmt.Errorf("read setup: %w", err)
	}
	return s, true, nil
}

func (Store) Begin(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, `INSERT INTO core.setup_state DEFAULT VALUES ON CONFLICT DO NOTHING`); err != nil {
		return fmt.Errorf("begin setup: %w", err)
	}
	return nil
}

func (Store) Sections(ctx context.Context, tx pgx.Tx) (map[string]SectionState, error) {
	rows, err := tx.Query(ctx, `
		SELECT section, status, assignee_id::text, to_char(due_on, 'YYYY-MM-DD'), note, updated_at FROM core.setup_sections`)
	if err != nil {
		return nil, fmt.Errorf("read sections: %w", err)
	}
	out := map[string]SectionState{}
	var s SectionState
	_, err = pgx.ForEachRow(rows, []any{&s.Key, &s.Status, &s.AssigneeID, &s.DueOn, &s.Note, &s.UpdatedAt}, func() error {
		out[s.Key] = s
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read sections: %w", err)
	}
	return out, nil
}

func (Store) PutSection(ctx context.Context, tx pgx.Tx, s SectionState, by string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO core.setup_sections (section, status, assignee_id, due_on, note, updated_by)
		VALUES ($1, $2, $3::uuid, $4::date, $5, $6::uuid)
		ON CONFLICT (tenant_id, section) DO UPDATE SET status = EXCLUDED.status, assignee_id = EXCLUDED.assignee_id,
			due_on = EXCLUDED.due_on, note = EXCLUDED.note, updated_by = EXCLUDED.updated_by, updated_at = now()`,
		s.Key, s.Status, s.AssigneeID, s.DueOn, s.Note, by)
	if err != nil {
		return fmt.Errorf("save section: %w", err)
	}
	return nil
}

func (Store) DueSoon(ctx context.Context, tx pgx.Tx, from, to string) ([]SectionState, error) {
	rows, err := tx.Query(ctx, `
		SELECT section, assignee_id::text, to_char(due_on, 'YYYY-MM-DD') FROM core.setup_sections
		WHERE assignee_id IS NOT NULL AND status <> 'approved' AND due_on BETWEEN $1::date AND $2::date
			AND reminded_on IS DISTINCT FROM due_on
		ORDER BY section`, from, to)
	if err != nil {
		return nil, fmt.Errorf("read due sections: %w", err)
	}
	out, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (SectionState, error) {
		var s SectionState
		err := row.Scan(&s.Key, &s.AssigneeID, &s.DueOn)
		return s, err
	})
	if err != nil {
		return nil, fmt.Errorf("read due sections: %w", err)
	}
	return out, nil
}

func (Store) MarkReminded(ctx context.Context, tx pgx.Tx, key, due string) error {
	if _, err := tx.Exec(ctx, `UPDATE core.setup_sections SET reminded_on = $2::date WHERE section = $1`, key, due); err != nil {
		return fmt.Errorf("mark reminded: %w", err)
	}
	return nil
}

func (Store) ClearAssignees(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, `UPDATE core.setup_sections SET assignee_id = NULL, due_on = NULL`); err != nil {
		return fmt.Errorf("clear assignees: %w", err)
	}
	return nil
}

func (Store) Legal(ctx context.Context, tx pgx.Tx) (map[string]LegalAnswer, error) {
	rows, err := tx.Query(ctx, `
		SELECT question, answer, basis, answered_by::text, answered_at, approved_by::text, approved_at FROM core.setup_legal`)
	if err != nil {
		return nil, fmt.Errorf("read legal: %w", err)
	}
	out := map[string]LegalAnswer{}
	var a LegalAnswer
	_, err = pgx.ForEachRow(rows, []any{&a.Question, &a.Answer, &a.Basis, &a.AnsweredBy, &a.AnsweredAt, &a.ApprovedBy, &a.ApprovedAt}, func() error {
		out[a.Question] = a
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read legal: %w", err)
	}
	return out, nil
}

func (Store) PutAnswer(ctx context.Context, tx pgx.Tx, question, answer, basis, by string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO core.setup_legal (question, answer, basis, answered_by) VALUES ($1, $2, $3, $4::uuid)
		ON CONFLICT (tenant_id, question) DO UPDATE SET answer = EXCLUDED.answer, basis = EXCLUDED.basis,
			answered_by = EXCLUDED.answered_by, answered_at = now(), approved_by = NULL, approved_at = NULL`,
		question, answer, basis, by)
	if err != nil {
		return fmt.Errorf("save answer: %w", err)
	}
	return nil
}

func (Store) ApproveAnswer(ctx context.Context, tx pgx.Tx, question, by string) (bool, error) {
	tag, err := tx.Exec(ctx, `
		UPDATE core.setup_legal SET approved_by = $2::uuid, approved_at = now() WHERE question = $1`, question, by)
	if err != nil {
		return false, fmt.Errorf("approve answer: %w", err)
	}
	return tag.RowsAffected() == 1, nil
}

func (Store) Unsign(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, `UPDATE core.setup_state SET signed_by = NULL, signed_at = NULL`); err != nil {
		return fmt.Errorf("unsign protocol: %w", err)
	}
	return nil
}

func (Store) Sign(ctx context.Context, tx pgx.Tx, by string) (int, error) {
	var version int
	err := tx.QueryRow(ctx, `
		UPDATE core.setup_state SET signed_by = $1::uuid, signed_at = now(), protocol_version = protocol_version + 1
		RETURNING protocol_version`, by).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("sign protocol: %w", err)
	}
	return version, nil
}

func (Store) Launch(ctx context.Context, tx pgx.Tx, by string, pilot []string) error {
	if pilot == nil {
		pilot = []string{}
	}
	_, err := tx.Exec(ctx, `UPDATE core.setup_state SET launched_by = $1::uuid, launched_at = now(), pilot_courses = $2::uuid[]`, by, pilot)
	if err != nil {
		return fmt.Errorf("launch: %w", err)
	}
	return nil
}

func (Store) Expand(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, `UPDATE core.setup_state SET pilot_courses = '{}'`); err != nil {
		return fmt.Errorf("expand pilot: %w", err)
	}
	return nil
}

func (Store) Hardware(ctx context.Context, tx pgx.Tx) (HardwareReport, error) {
	r := HardwareReport{Checks: []module.HardwareCheck{}}
	var raw []byte
	if err := tx.QueryRow(ctx, `SELECT hardware_checks, hardware_at FROM core.setup_state`).Scan(&raw, &r.CheckedAt); err != nil {
		return HardwareReport{}, fmt.Errorf("read hardware checks: %w", err)
	}
	if raw != nil {
		if err := json.Unmarshal(raw, &r.Checks); err != nil {
			return HardwareReport{}, fmt.Errorf("decode hardware checks: %w", err)
		}
	}
	return r, nil
}

func (Store) PutHardware(ctx context.Context, tx pgx.Tx, checks []module.HardwareCheck) (HardwareReport, error) {
	raw, err := json.Marshal(checks)
	if err != nil {
		return HardwareReport{}, err
	}
	r := HardwareReport{Checks: checks}
	if err := tx.QueryRow(ctx, `UPDATE core.setup_state SET hardware_checks = $1, hardware_at = now() RETURNING hardware_at`, raw).Scan(&r.CheckedAt); err != nil {
		return HardwareReport{}, fmt.Errorf("save hardware checks: %w", err)
	}
	return r, nil
}
