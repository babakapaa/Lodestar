package team_test

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/team"
)

func ids(list ...auth.Principal) []string {
	out := make([]string, 0, len(list))
	for _, p := range list {
		out = append(out, p.UserID)
	}
	return out
}

func TestStudentsGatherTeamsThemselves(t *testing.T) {
	w := newWorld(t)
	teacher, aigerim, erlan, dana, timur := w.user("dana-t"), w.user("aigerim"), w.user("erlan"), w.user("dana"), w.user("timur")
	outsider := w.user("outsider")
	later := time.Now().Add(72 * time.Hour)
	w.scopes.set(module.TeamScope{
		ID: "a1", Title: "Курсовая", Context: "Базы данных", Leads: ids(teacher), Members: ids(aigerim, erlan, dana, timur),
		MinSize: 2, MaxSize: 3, SelfFormed: true, FreeLeaveUntil: &later,
	})
	if _, err := w.teams.View(w.ctx, outsider, "study.assignment", "a1"); !errors.Is(err, team.ErrScopeNotFound) {
		t.Fatalf("outsider view: %v", err)
	}

	alpha, err := w.teams.Create(w.ctx, aigerim, "study.assignment", "a1", team.TeamInput{Name: " Альфа "})
	if err != nil || alpha.Name != "Альфа" || len(alpha.Members) != 1 || alpha.InviteCode == nil || !alpha.CanLeaveFreely {
		t.Fatalf("create = %+v, %v", alpha, err)
	}
	if _, err := w.teams.Create(w.ctx, aigerim, "study.assignment", "a1", team.TeamInput{Name: "Вторая"}); !errors.Is(err, team.ErrInTeam) {
		t.Fatalf("second team: %v", err)
	}

	if _, err := w.teams.Invite(w.ctx, aigerim, alpha.ID, team.InviteInput{UserID: erlan.UserID}); err != nil {
		t.Fatalf("invite: %v", err)
	}
	if got := w.sent.of(team.NotifyInvited); len(got) != 1 || got[0].UserIDs[0] != erlan.UserID || got[0].Params["inviter"] != "aigerim" {
		t.Fatalf("invited = %+v", got)
	}
	cards, err := w.teams.Invitations(w.ctx, erlan)
	if err != nil || len(cards) != 1 || cards[0].InvitedBy == nil || cards[0].Scope.FreeLeaveUntil == nil {
		t.Fatalf("invitations = %+v, %v", cards, err)
	}
	if _, err := w.teams.Accept(w.ctx, dana, alpha.ID); !errors.Is(err, team.ErrNotInvited) {
		t.Fatalf("accept without invite: %v", err)
	}
	if alpha, err = w.teams.Accept(w.ctx, erlan, alpha.ID); err != nil || len(alpha.Members) != 2 {
		t.Fatalf("accept = %+v, %v", alpha, err)
	}

	preview, err := w.teams.Preview(w.ctx, dana, *alpha.InviteCode)
	if err != nil || preview.Team.Name != "Альфа" || preview.Scope.MaxSize != 3 {
		t.Fatalf("preview = %+v, %v", preview, err)
	}
	if _, err := w.teams.Preview(w.ctx, outsider, *alpha.InviteCode); !errors.Is(err, team.ErrTeamNotFound) {
		t.Fatalf("outsider preview: %v", err)
	}
	if _, err := w.teams.Join(w.ctx, dana, *alpha.InviteCode); err != nil {
		t.Fatalf("join: %v", err)
	}
	if _, err := w.teams.Join(w.ctx, timur, *alpha.InviteCode); !errors.Is(err, team.ErrTeamFull) {
		t.Fatalf("join full: %v", err)
	}
	view, err := w.teams.Seek(w.ctx, timur, "study.assignment", "a1", team.SeekingInput{Seeking: true, Note: "Пишу на Go"})
	if err != nil || !view.Seeking || len(view.Seekers) != 1 || view.Teams != nil {
		t.Fatalf("seek = %+v, %v", view, err)
	}

	if alpha, err = w.teams.Label(w.ctx, erlan, alpha.ID, team.LabelInput{Label: "код"}); err != nil {
		t.Fatalf("label: %v", err)
	}
	if i := slices.IndexFunc(alpha.Members, func(m team.Member) bool { return m.UserID == erlan.UserID }); alpha.Members[i].Label != "код" {
		t.Fatalf("label not kept: %+v", alpha.Members)
	}
	left, err := w.teams.Leave(w.ctx, dana, alpha.ID, team.LeaveInput{})
	if err != nil || !left.Left {
		t.Fatalf("leave = %+v, %v", left, err)
	}
	if r := w.roster(alpha.ID); !r.Group || len(r.Members) != 2 || r.Archived || r.Title != "Альфа" {
		t.Fatalf("roster = %+v", r)
	}

	view, err = w.teams.View(w.ctx, teacher, "study.assignment", "a1")
	if err != nil || view.Role != team.RoleLead || len(view.Teams) != 1 || len(view.Unplaced) != 2 {
		t.Fatalf("teacher view = %+v, %v", view, err)
	}
	view, err = w.teams.Distribute(w.ctx, teacher, "study.assignment", "a1")
	if err != nil || len(view.Teams) != 2 || len(view.Unplaced) != 0 || len(view.Seekers) != 0 {
		t.Fatalf("distribute = %+v, %v", view, err)
	}
	if got := w.sent.of(team.NotifyChanged); len(got) != 1 || len(got[0].UserIDs) != 2 {
		t.Fatalf("changed = %+v", got)
	}
	mine, err := w.teams.Mine(w.ctx, timur)
	if err != nil || len(mine) != 1 || mine[0].Scope.Title != "Курсовая" {
		t.Fatalf("mine = %+v, %v", mine, err)
	}
}

func TestTeacherFormsTeamsAndDecidesOnLeaving(t *testing.T) {
	w := newWorld(t)
	teacher := w.user("teacher")
	s1, s2, s3, s4 := w.user("s1"), w.user("s2"), w.user("s3"), w.user("s4")
	earlier := time.Now().Add(-time.Hour)
	scope := module.TeamScope{
		ID: "b1", Title: "Проект", Leads: ids(teacher), Members: ids(s1, s2, s3, s4),
		MinSize: 2, MaxSize: 4, FreeLeaveUntil: &earlier,
	}
	w.scopes.set(scope)
	if _, err := w.teams.Create(w.ctx, s1, "study.assignment", "b1", team.TeamInput{Name: "Своя"}); !errors.Is(err, team.ErrNotSelfFormed) {
		t.Fatalf("student creates: %v", err)
	}
	first, err := w.teams.Create(w.ctx, teacher, "study.assignment", "b1", team.TeamInput{Name: "Первая", Members: ids(s1, s2, s3)})
	if err != nil || len(first.Members) != 3 || !first.Manage {
		t.Fatalf("teacher creates = %+v, %v", first, err)
	}
	if _, err := w.teams.Invite(w.ctx, s1, first.ID, team.InviteInput{UserID: s4.UserID}); !errors.Is(err, team.ErrNotSelfFormed) {
		t.Fatalf("student invites: %v", err)
	}
	if d, err := w.teams.Get(w.ctx, s1, first.ID); err != nil || d.InviteCode != nil || d.CanLeaveFreely {
		t.Fatalf("member details = %+v, %v", d, err)
	}

	asked, err := w.teams.Leave(w.ctx, s2, first.ID, team.LeaveInput{Reason: "Не сходимся по графику"})
	if err != nil || asked.Left || asked.Request == nil || asked.Request.Reason != "Не сходимся по графику" {
		t.Fatalf("leave request = %+v, %v", asked, err)
	}
	if _, err := w.teams.Leave(w.ctx, s2, first.ID, team.LeaveInput{}); !errors.Is(err, team.ErrRequestPending) {
		t.Fatalf("second request: %v", err)
	}
	if got := w.sent.of(team.NotifyLeaveRequested); len(got) != 1 || got[0].UserIDs[0] != teacher.UserID {
		t.Fatalf("requested = %+v", got)
	}
	if _, err := w.teams.Decide(w.ctx, s1, asked.Request.ID, team.DecisionInput{Approve: true}); !errors.Is(err, team.ErrRequestNotFound) {
		t.Fatalf("student decides: %v", err)
	}
	if d, err := w.teams.Decide(w.ctx, teacher, asked.Request.ID, team.DecisionInput{}); err != nil || len(d.Members) != 3 {
		t.Fatalf("reject = %+v, %v", d, err)
	}
	asked, err = w.teams.Leave(w.ctx, s2, first.ID, team.LeaveInput{})
	if err != nil {
		t.Fatalf("ask again: %v", err)
	}
	if d, err := w.teams.Decide(w.ctx, teacher, asked.Request.ID, team.DecisionInput{Approve: true}); err != nil || len(d.Members) != 2 {
		t.Fatalf("approve = %+v, %v", d, err)
	}
	if len(w.sent.of(team.NotifyLeaveRejected)) != 1 || len(w.sent.of(team.NotifyLeaveApproved)) != 1 {
		t.Fatal("decisions not notified")
	}

	second, err := w.teams.Create(w.ctx, teacher, "study.assignment", "b1", team.TeamInput{Name: "Вторая", Members: ids(s2, s4)})
	if err != nil {
		t.Fatalf("second: %v", err)
	}
	if second, err = w.teams.Add(w.ctx, teacher, second.ID, s3.UserID); err != nil || len(second.Members) != 3 {
		t.Fatalf("move = %+v, %v", second, err)
	}
	if first, err = w.teams.Remove(w.ctx, teacher, first.ID, s1.UserID, team.RemovalInput{Solo: true}); err != nil || first.ArchivedAt == nil {
		t.Fatalf("separate = %+v, %v", first, err)
	}
	if _, err := w.teams.Get(w.ctx, s1, first.ID); !errors.Is(err, team.ErrTeamNotFound) {
		t.Fatalf("separated reads team: %v", err)
	}
	if _, err := w.teams.Seek(w.ctx, s1, "study.assignment", "b1", team.SeekingInput{Seeking: true}); !errors.Is(err, team.ErrSeparated) {
		t.Fatalf("separated seeks: %v", err)
	}
	var events []team.HistoryEvent
	for _, e := range first.History {
		events = append(events, e.Event)
	}
	want := []team.HistoryEvent{team.EventCreated, team.EventAdded, team.EventAdded, team.EventAdded, team.EventRemoved, team.EventMovedOut, team.EventSeparated}
	if !slices.Equal(events, want) {
		t.Fatalf("history = %v", events)
	}
	view, err := w.teams.View(w.ctx, teacher, "study.assignment", "b1")
	if err != nil || len(view.Teams) != 1 || len(view.Solos) != 1 || len(view.Unplaced) != 0 {
		t.Fatalf("view = %+v, %v", view, err)
	}

	scope.Archived = true
	w.scopes.set(scope)
	if _, err := w.teams.Remove(w.ctx, teacher, second.ID, s4.UserID, team.RemovalInput{}); err != nil {
		t.Fatalf("remove from archived scope: %v", err)
	}
	if _, err := w.teams.Add(w.ctx, teacher, second.ID, s4.UserID); !errors.Is(err, team.ErrScopeArchived) {
		t.Fatalf("add in archived scope: %v", err)
	}
	if r := w.roster(second.ID); !r.Archived {
		t.Fatalf("roster not archived: %+v", r)
	}
}
