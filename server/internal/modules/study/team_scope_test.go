package study_test

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestTeamAssignmentsGiveTeamsTheirRules(t *testing.T) {
	d := newDefenseCase(t, study.DefenseShort)
	w := d.w
	scopes := study.Module{}.TeamScopeSources(module.Kernel{})[0]
	load := func(id string) (module.TeamScope, error) {
		var s module.TeamScope
		err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
			var err error
			s, err = scopes.TeamScope(w.ctx, tx, id)
			return err
		})
		return s, err
	}
	if _, err := load(d.work.ID); !errors.Is(err, module.ErrTeamScopeNotFound) {
		t.Fatalf("solo assignment: %v", err)
	}
	if _, err := w.work.Create(w.ctx, d.teacher, d.course.ID, study.AssignmentInput{
		Title: "Команды", Team: &study.TeamRules{MinSize: 1, MaxSize: 3},
	}); !errors.Is(err, study.ErrInvalidAssignment) {
		t.Fatalf("team of one: %v", err)
	}
	until := time.Now().Add(48 * time.Hour).Truncate(time.Second)
	rules := study.TeamRules{MinSize: 2, MaxSize: 4, SelfFormed: true, FreeLeaveUntil: &until}
	draft, err := w.work.Create(w.ctx, d.teacher, d.course.ID, study.AssignmentInput{Title: "Проект", Team: &rules})
	if err != nil || draft.Team == nil || draft.Team.MaxSize != 4 {
		t.Fatalf("create = %+v, %v", draft, err)
	}
	if _, err := load(draft.ID); !errors.Is(err, module.ErrTeamScopeNotFound) {
		t.Fatalf("draft scope: %v", err)
	}
	published, err := w.work.Update(w.ctx, d.teacher, draft.ID, study.AssignmentInput{Title: "Проект", Team: &rules, Published: true})
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	s, err := load(published.ID)
	if err != nil || s.Title != "Проект" || s.Context != "Базы данных" || !s.SelfFormed || !s.FreeLeaveUntil.Equal(until) {
		t.Fatalf("scope = %+v, %v", s, err)
	}
	if !slices.Equal(s.Leads, []string{d.teacher.UserID}) || len(s.Members) != 2 || !slices.Contains(s.Members, d.first.UserID) {
		t.Fatalf("people = %v / %v", s.Leads, s.Members)
	}
}
