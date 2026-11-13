package league

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const (
	TaskScopeKind    = "league.task"
	similarThreshold = 0.9
)

type TaskScopes struct {
	service Service
}

func (TaskScopes) Kind() string {
	return TaskScopeKind
}

func (s TaskScopes) TeamScope(ctx context.Context, tx pgx.Tx, id string) (module.TeamScope, error) {
	t, r, err := s.load(ctx, tx, id)
	if err != nil {
		return module.TeamScope{}, err
	}
	return module.TeamScope{
		Kind: TaskScopeKind, ID: t.ID, Title: t.Title, Context: r.Title, Link: taskLink(t.ID),
		Leads: r.Leads, Members: r.Members, MinSize: t.MinSize, MaxSize: t.MaxSize, SelfFormed: t.SelfFormed,
		DueAt: t.DueAt, Archived: t.ClosedAt != nil || r.Archived,
	}, nil
}

func (s TaskScopes) PrepareTeamWork(ctx context.Context, tx pgx.Tx, id string, p auth.Principal, work json.RawMessage) (json.RawMessage, error) {
	t, r, err := s.load(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(r.Members, p.UserID) {
		return nil, module.ErrTeamScopeNotFound
	}
	if t.ClosedAt != nil || r.Archived {
		return nil, module.TeamWorkProblem{Status: http.StatusConflict, Code: CodeTaskClosed}
	}
	var in SolutionWork
	if err := json.Unmarshal(work, &in); err != nil {
		return nil, module.TeamWorkProblem{Status: http.StatusBadRequest, Code: CodeInvalidSolution}
	}
	if in, err = in.normalized(); err != nil {
		return nil, module.TeamWorkProblem{Status: http.StatusBadRequest, Code: CodeInvalidSolution}
	}
	return json.Marshal(in)
}

func (s TaskScopes) DeliverTeamWork(ctx context.Context, tx pgx.Tx, d module.TeamDelivery) error {
	t, r, err := s.load(ctx, tx, d.ScopeID)
	if err != nil {
		return err
	}
	var work SolutionWork
	if err := json.Unmarshal(d.Work, &work); err != nil {
		return err
	}
	sol := Solution{TaskID: t.ID, TeamID: d.TeamID, TeamName: d.TeamName, Work: work, Votes: []SolutionVote{}}
	if d.ProposedBy != "" {
		sol.ProposedBy = &d.ProposedBy
	}
	for _, v := range d.Votes {
		sol.Members = append(sol.Members, v.UserID)
		sol.Votes = append(sol.Votes, SolutionVote{UserID: v.UserID, Decision: v.Decision, Comment: v.Comment})
	}
	saved, created, err := s.service.solutions.Insert(ctx, tx, d.ProposalID, sol)
	if err != nil || !created {
		return err
	}
	if err := s.compare(ctx, tx, saved); err != nil {
		return err
	}
	if len(r.Leads) == 0 {
		return nil
	}
	return s.service.notifier.Notify(ctx, tx, d.TenantID, module.Notification{
		Type: NotifySolution, UserIDs: r.Leads, Params: map[string]any{"task": t.Title, "team": d.TeamName}, Link: taskLink(t.ID),
	})
}

func (s TaskScopes) compare(ctx context.Context, tx pgx.Tx, sol Solution) error {
	if s.service.embedder == nil {
		return nil
	}
	text := strings.TrimSpace(sol.Work.Text + "\n" + strings.Join(sol.Work.Links, "\n"))
	vectors, err := s.service.embedder.Embed(ctx, []string{text})
	if err != nil || len(vectors) != 1 {
		slog.Warn("embed league solution", "error", err)
		return nil
	}
	others, err := s.service.solutions.Embeddings(ctx, tx, sol.TeamID)
	if err != nil {
		return err
	}
	var best *string
	var score *float64
	for id, other := range others {
		similarity := cosine(vectors[0], other)
		if similarity >= similarThreshold && (score == nil || similarity > *score) {
			best, score = &id, &similarity
		}
	}
	return s.service.solutions.Embed(ctx, tx, sol.ID, vectors[0], best, score)
}

func (s TaskScopes) load(ctx context.Context, tx pgx.Tx, id string) (Task, module.Roster, error) {
	t, err := s.service.tasks.ByID(ctx, tx, id)
	if errors.Is(err, ErrTaskNotFound) {
		return t, module.Roster{}, module.ErrTeamScopeNotFound
	}
	if err != nil {
		return t, module.Roster{}, err
	}
	l, err := s.service.leagues.ByID(ctx, tx, t.LeagueID)
	if err != nil {
		return t, module.Roster{}, err
	}
	r, _, err := s.service.arena(ctx, tx, l.OriginKind, l.OriginID)
	if errors.Is(err, ErrLeagueNotFound) {
		return t, r, module.ErrTeamScopeNotFound
	}
	return t, r, err
}

func cosine(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, na, nb float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		na += float64(a[i]) * float64(a[i])
		nb += float64(b[i]) * float64(b[i])
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}
