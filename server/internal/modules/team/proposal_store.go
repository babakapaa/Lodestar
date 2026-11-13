package team

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const proposalColumns = `p.id::text, p.team_id::text, u.id::text, u.display_name, u.email, p.work, p.note, p.status, p.silent_at, p.due_at, p.created_at, p.sent_at`

type ProposalStore struct{}

func (ProposalStore) Insert(ctx context.Context, tx pgx.Tx, teamID, proposedBy string, work json.RawMessage, note string, silent time.Time, due *time.Time) (string, error) {
	if _, err := tx.Exec(ctx, `UPDATE team.proposals SET status = 'replaced' WHERE team_id::text = $1 AND status = 'open'`, teamID); err != nil {
		return "", fmt.Errorf("replace proposal: %w", err)
	}
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO team.proposals (team_id, proposed_by, work, note, silent_at, due_at)
		VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6) RETURNING id::text`, teamID, proposedBy, work, note, silent, due).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("propose: %w", err)
	}
	return id, nil
}

func (ProposalStore) Vote(ctx context.Context, tx pgx.Tx, proposalID, userID string, agreed bool, comment string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO team.votes (proposal_id, user_id, agreed, comment) VALUES ($1::uuid, $2::uuid, $3, $4)
		ON CONFLICT (proposal_id, user_id) DO UPDATE SET agreed = excluded.agreed, comment = excluded.comment, at = now()`,
		proposalID, userID, agreed, comment)
	if err != nil {
		return fmt.Errorf("vote: %w", err)
	}
	return nil
}

func (ProposalStore) Sent(ctx context.Context, tx pgx.Tx, id string) error {
	if _, err := tx.Exec(ctx, `UPDATE team.proposals SET status = 'sent', sent_at = now() WHERE id::text = $1`, id); err != nil {
		return fmt.Errorf("send proposal: %w", err)
	}
	return nil
}

func (s ProposalStore) ByID(ctx context.Context, tx pgx.Tx, id string, lock bool) (Proposal, error) {
	query := `SELECT ` + proposalColumns + ` FROM team.proposals p LEFT JOIN core.users u ON u.id = p.proposed_by WHERE p.id::text = $1`
	if lock {
		query += ` FOR UPDATE OF p`
	}
	list, err := s.many(ctx, tx, query, id)
	if err != nil {
		return Proposal{}, err
	}
	if len(list) == 0 {
		return Proposal{}, ErrProposalNotFound
	}
	return list[0], nil
}

func (s ProposalStore) Latest(ctx context.Context, tx pgx.Tx, teamID string) (*Proposal, error) {
	list, err := s.many(ctx, tx, `
		SELECT `+proposalColumns+` FROM team.proposals p LEFT JOIN core.users u ON u.id = p.proposed_by
		WHERE p.team_id::text = $1 AND p.status <> 'replaced'
		ORDER BY p.status = 'open' DESC, p.created_at DESC LIMIT 1`, teamID)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return &list[0], nil
}

func (ProposalStore) Due(ctx context.Context, tx pgx.Tx, now time.Time) ([]string, error) {
	rows, err := tx.Query(ctx, `
		SELECT id::text FROM team.proposals
		WHERE status = 'open' AND (silent_at <= $1 OR due_at <= $1) ORDER BY created_at`, now)
	if err != nil {
		return nil, fmt.Errorf("due proposals: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (ProposalStore) Votes(ctx context.Context, tx pgx.Tx, proposalID string, members []Member) ([]Vote, error) {
	rows, err := tx.Query(ctx, `SELECT user_id::text, agreed, comment, at FROM team.votes WHERE proposal_id::text = $1`, proposalID)
	if err != nil {
		return nil, fmt.Errorf("votes: %w", err)
	}
	type cast struct {
		agreed  bool
		comment string
		at      time.Time
	}
	given := map[string]cast{}
	defer rows.Close()
	for rows.Next() {
		var userID string
		var c cast
		if err := rows.Scan(&userID, &c.agreed, &c.comment, &c.at); err != nil {
			return nil, fmt.Errorf("votes: %w", err)
		}
		given[userID] = c
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("votes: %w", err)
	}
	votes := make([]Vote, 0, len(members))
	for _, m := range members {
		v := Vote{Person: m.Person, Decision: VoteWaiting}
		if c, ok := given[m.UserID]; ok {
			at := c.at
			v.Decision, v.Comment, v.At = module.VoteObjected, c.comment, &at
			if c.agreed {
				v.Decision = module.VoteAgreed
			}
		}
		votes = append(votes, v)
	}
	return votes, nil
}

func (ProposalStore) many(ctx context.Context, tx pgx.Tx, query string, args ...any) ([]Proposal, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("proposals: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Proposal, error) {
		var p Proposal
		var userID, name, email *string
		err := row.Scan(&p.ID, &p.TeamID, &userID, &name, &email, &p.Work, &p.Note, &p.Status, &p.SilentAt, &p.DueAt, &p.CreatedAt, &p.SentAt)
		p.ProposedBy = person(userID, name, email)
		return p, err
	})
}
