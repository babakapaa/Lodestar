package study

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s AssignmentService) teamless(ctx context.Context, tx pgx.Tx, a Assignment, userID string) error {
	if a.Team == nil || s.teams == nil {
		return nil
	}
	placement, err := s.teams.Placement(ctx, tx, TeamScopeAssignment, a.ID, userID)
	if err != nil {
		return err
	}
	if placement.TeamID != "" {
		return ErrTeamWork
	}
	return nil
}

func (s AssignmentService) acceptable(ctx context.Context, tx pgx.Tx, p auth.Principal, a Assignment, c Course, in SubmissionInput) error {
	if c.ArchivedAt != nil {
		return ErrArchived
	}
	if a.SearchMap != nil && (in.SearchMap == nil || len(in.SearchMap.Missing(*a.SearchMap)) > 0) {
		return ErrSearchMapIncomplete
	}
	return s.attachments.Check(ctx, tx, p.UserID, in.Files)
}

func (s AssignmentService) PrepareTeamWork(ctx context.Context, tx pgx.Tx, id string, p auth.Principal, work json.RawMessage) (json.RawMessage, error) {
	var in SubmissionInput
	if err := json.Unmarshal(work, &in); err != nil {
		return nil, problem(ErrInvalidSubmission)
	}
	in, err := in.content()
	if err != nil {
		return nil, problem(err)
	}
	a, c, st, err := s.visible(ctx, tx, p, id, false)
	if err != nil {
		return nil, problem(err)
	}
	if st.Role != RoleStudent || a.Team == nil {
		return nil, module.ErrTeamScopeNotFound
	}
	if err := s.acceptable(ctx, tx, p, a, c, in); err != nil {
		return nil, problem(err)
	}
	in.ClientID = ""
	snapshot, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("encode team work: %w", err)
	}
	return snapshot, nil
}

func (s AssignmentService) DeliverTeamWork(ctx context.Context, tx pgx.Tx, d module.TeamDelivery) error {
	var in SubmissionInput
	if err := json.Unmarshal(d.Work, &in); err != nil {
		return fmt.Errorf("decode team work: %w", err)
	}
	in.ClientID = d.ProposalID
	a, err := s.assignments.ByID(ctx, tx, d.ScopeID, true)
	if err != nil {
		return err
	}
	c, err := CourseStore{}.ByID(ctx, tx, a.CourseID)
	if err != nil {
		return err
	}
	receipt, err := s.receipt(ctx, tx, d)
	if err != nil {
		return err
	}
	for _, v := range d.Votes {
		if _, err := s.submissions.ByClient(ctx, tx, a.ID, v.UserID, in.ClientID); err == nil {
			continue
		} else if !errors.Is(err, ErrAssignmentNotFound) {
			return err
		}
		help, err := s.helpFor(ctx, tx, a, v.UserID)
		if err != nil {
			return err
		}
		sub, err := s.submissions.Insert(ctx, tx, a.ID, v.UserID, in, &receipt, help)
		if err != nil {
			return err
		}
		if err := s.submissions.SetFiles(ctx, tx, sub.ID, in.Files); err != nil {
			return err
		}
	}
	return s.tellTeachers(ctx, tx, d.TenantID, c, a, d.TeamName)
}

func (s AssignmentService) receipt(ctx context.Context, tx pgx.Tx, d module.TeamDelivery) (TeamReceipt, error) {
	ids := make([]string, 0, len(d.Votes))
	for _, v := range d.Votes {
		ids = append(ids, v.UserID)
	}
	people, err := s.people.ByIDs(ctx, tx, ids)
	if err != nil {
		return TeamReceipt{}, err
	}
	byID := make(map[string]identity.Person, len(people))
	for _, person := range people {
		byID[person.UserID] = person
	}
	r := TeamReceipt{TeamID: d.TeamID, TeamName: d.TeamName, ProposalID: d.ProposalID, Votes: make([]ReceiptVote, 0, len(d.Votes))}
	if author, ok := byID[d.ProposedBy]; ok {
		r.ProposedBy = &author
	}
	for _, v := range d.Votes {
		r.Votes = append(r.Votes, ReceiptVote{Person: byID[v.UserID], Decision: v.Decision, Comment: v.Comment})
	}
	return r, nil
}
