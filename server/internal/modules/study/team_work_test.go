package study_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

type placedIn string

func (p placedIn) Placement(context.Context, pgx.Tx, string, string, string) (module.TeamPlacement, error) {
	return module.TeamPlacement{TeamID: string(p)}, nil
}

func TestTeamWorkIsCheckedAndHandedInForEveryMember(t *testing.T) {
	d := newDefenseCase(t, study.DefenseShort)
	w := d.w
	rules := study.TeamRules{MinSize: 2, MaxSize: 3}
	a, err := w.work.Create(w.ctx, d.teacher, d.course.ID, study.AssignmentInput{
		Title: "Проект", Team: &rules, SearchMap: &study.MapRules{MinSources: 1}, Published: true,
	})
	if err != nil {
		t.Fatalf("assignment: %v", err)
	}
	scopes := study.Module{}.TeamScopeSources(w.kernel())[0]
	prepare := func(p auth.Principal, work string) (json.RawMessage, error) {
		var out json.RawMessage
		err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
			var err error
			out, err = scopes.PrepareTeamWork(w.ctx, tx, a.ID, p, json.RawMessage(work))
			return err
		})
		return out, err
	}
	var problem module.TeamWorkProblem
	if _, err := prepare(d.first, `{"body":"Без карты"}`); !errors.As(err, &problem) || problem.Code != study.CodeMapIncomplete {
		t.Fatalf("no map: %v", err)
	}
	if _, err := prepare(d.teacher, `{"body":"x"}`); !errors.Is(err, module.ErrTeamScopeNotFound) {
		t.Fatalf("teacher prepares: %v", err)
	}
	m := `{"used":[{"reference":"Codd E. F.","kind":"article","found":"eLibrary","why":"Основа теории"}],"rejected":[],"path":"Нашли первоисточник"}`
	snapshot, err := prepare(d.first, `{"clientId":"x","body":"  Отчёт команды  ","searchMap":`+m+`}`)
	if err != nil {
		t.Fatalf("prepare: %v", err)
	}
	delivery := module.TeamDelivery{
		TenantID: w.tenant, ScopeID: a.ID, TeamID: "4c2b1f0e-7a8d-4e3b-9c51-2d6f0a1b3e77", TeamName: "Альфа",
		ProposalID: "9a1b2c3d-4e5f-4a6b-8c7d-0e1f2a3b4c5d", ProposedBy: d.first.UserID, Work: snapshot,
		Votes: []module.TeamVote{
			{UserID: d.first.UserID, Decision: module.VoteAgreed},
			{UserID: d.second.UserID, Decision: module.VoteObjected, Comment: "Мало источников"},
		},
	}
	for range 2 {
		err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
			return scopes.DeliverTeamWork(w.ctx, tx, delivery)
		})
		if err != nil {
			t.Fatalf("deliver: %v", err)
		}
	}
	for _, student := range []auth.Principal{d.first, d.second} {
		attempts, err := w.work.Attempts(w.ctx, d.teacher, a.ID, student.UserID)
		if err != nil || len(attempts) != 1 {
			t.Fatalf("attempts = %+v, %v", attempts, err)
		}
		sub := attempts[0]
		if sub.Body != "Отчёт команды" || sub.SearchMap == nil || sub.Team == nil || sub.Team.TeamName != "Альфа" {
			t.Fatalf("submission = %+v", sub)
		}
		if sub.Team.ProposedBy == nil || sub.Team.ProposedBy.UserID != d.first.UserID || sub.Team.Votes[1].Comment != "Мало источников" {
			t.Fatalf("receipt = %+v", sub.Team)
		}
	}

	in := study.SubmissionInput{ClientID: "1f0e4c2b-7a8d-4e3b-9c51-2d6f0a1b3e77", Body: "Сам", SearchMap: &study.SearchMap{
		Used: []study.MapSource{{Reference: "Codd E. F.", Kind: "article", Found: "eLibrary", Why: "Основа теории"}}, Path: "Нашёл сам",
	}}
	if _, _, err := w.work.WithTeams(placedIn("t1")).Submit(w.ctx, d.second, a.ID, in); !errors.Is(err, study.ErrTeamWork) {
		t.Fatalf("member submits alone: %v", err)
	}
	if _, _, err := w.work.WithTeams(placedIn("")).Submit(w.ctx, d.second, a.ID, in); err != nil {
		t.Fatalf("student without team: %v", err)
	}
}
