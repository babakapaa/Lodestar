package team_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/team"
)

func work(body string) team.ProposalInput {
	return team.ProposalInput{Work: json.RawMessage(`{"body":"` + body + `"}`)}
}

func TestTeamHandsInWorkWhenEveryoneAgrees(t *testing.T) {
	w := newWorld(t)
	teacher, a, b, c := w.user("teacher"), w.user("aigerim"), w.user("erlan"), w.user("dana")
	due := time.Now().Add(7 * 24 * time.Hour)
	w.scopes.set(module.TeamScope{ID: "p1", Title: "Проект", Leads: ids(teacher), Members: ids(a, b, c), MinSize: 2, MaxSize: 3, DueAt: &due})
	d, err := w.teams.Create(w.ctx, teacher, "study.assignment", "p1", team.TeamInput{Name: "Альфа", Members: ids(a, b, c)})
	if err != nil {
		t.Fatalf("team: %v", err)
	}
	if _, err := w.teams.Propose(w.ctx, a, d.ID, work("incomplete")); err == nil {
		t.Fatal("scope rejection lost")
	} else if problem := (module.TeamWorkProblem{}); !errors.As(err, &problem) || problem.Code != "search_map_incomplete" {
		t.Fatalf("problem = %v", err)
	}
	if _, err := w.teams.Propose(w.ctx, teacher, d.ID, work("x")); !errors.Is(err, team.ErrTeamNotFound) {
		t.Fatalf("teacher proposes: %v", err)
	}
	d, err = w.teams.Propose(w.ctx, a, d.ID, work("Версия 1"))
	if err != nil || d.Proposal == nil || d.Proposal.Status != team.ProposalOpen || len(d.Proposal.Votes) != 3 {
		t.Fatalf("propose = %+v, %v", d.Proposal, err)
	}
	if got := w.sent.of(team.NotifyProposed); len(got) != 1 || len(got[0].UserIDs) != 2 {
		t.Fatalf("proposed = %+v", got)
	}
	d, err = w.teams.Propose(w.ctx, b, d.ID, work("Версия 2"))
	if err != nil || d.Proposal.ProposedBy.UserID != b.UserID {
		t.Fatalf("second proposal = %+v, %v", d.Proposal, err)
	}
	if _, err := w.teams.Vote(w.ctx, c, d.Proposal.ID, team.VoteInput{}); !errors.Is(err, team.ErrInvalidTeam) {
		t.Fatalf("objection without comment: %v", err)
	}
	if d, err = w.teams.Vote(w.ctx, c, d.Proposal.ID, team.VoteInput{Agree: true}); err != nil || d.Proposal.Status != team.ProposalOpen {
		t.Fatalf("vote = %+v, %v", d.Proposal, err)
	}
	if d, err = w.teams.Vote(w.ctx, a, d.Proposal.ID, team.VoteInput{Agree: true}); err != nil || d.Proposal.Status != team.ProposalSent {
		t.Fatalf("last vote = %+v, %v", d.Proposal, err)
	}
	got := w.scopes.deliveries()
	var delivered struct{ Body string }
	if len(got) != 1 || json.Unmarshal(got[0].Work, &delivered) != nil || delivered.Body != "Версия 2" || got[0].ProposedBy != b.UserID || len(got[0].Votes) != 3 {
		t.Fatalf("delivered = %+v", got)
	}
	for _, v := range got[0].Votes {
		if v.Decision != module.VoteAgreed {
			t.Fatalf("vote = %+v", v)
		}
	}
	if n := len(w.sent.of(team.NotifySubmitted)); n != 1 {
		t.Fatalf("submitted notified %d times", n)
	}
	if _, err := w.teams.Vote(w.ctx, c, d.Proposal.ID, team.VoteInput{Agree: true}); !errors.Is(err, team.ErrProposalNotFound) {
		t.Fatalf("vote after send: %v", err)
	}
	placement, err := w.placement("p1", c.UserID)
	if err != nil || placement.TeamID != d.ID {
		t.Fatalf("placement = %+v, %v", placement, err)
	}
	var refs []module.TeamRef
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		refs, err = team.Module{}.TeamDirectory(w.kernel).Teams(w.ctx, tx, c.UserID)
		return err
	})
	if err != nil || len(refs) != 1 || refs[0].ID != d.ID || len(refs[0].Members) != len(d.Members) {
		t.Fatalf("teams of member = %+v, %v", refs, err)
	}
}

func TestSilenceAgreesAndObjectionWaitsForTheDeadline(t *testing.T) {
	w := newWorld(t)
	teacher, a, b := w.user("teacher"), w.user("aigerim"), w.user("erlan")
	due := time.Now().Add(72 * time.Hour)
	w.scopes.set(module.TeamScope{ID: "p2", Title: "Проект", Leads: ids(teacher), Members: ids(a, b), MinSize: 2, MaxSize: 2, DueAt: &due})
	d, err := w.teams.Create(w.ctx, teacher, "study.assignment", "p2", team.TeamInput{Name: "Пара", Members: ids(a, b)})
	if err != nil {
		t.Fatalf("team: %v", err)
	}
	settler := team.NewWorkSettler(w.kernel)

	d, err = w.teams.Propose(w.ctx, a, d.ID, work("Черновик"))
	if err != nil {
		t.Fatalf("propose: %v", err)
	}
	if !d.Proposal.SilentAt.Before(time.Now().Add(25 * time.Hour)) {
		t.Fatalf("silent at = %v", d.Proposal.SilentAt)
	}
	if d, err = w.teams.Vote(w.ctx, b, d.Proposal.ID, team.VoteInput{Comment: "Нет выводов"}); err != nil {
		t.Fatalf("object: %v", err)
	}
	if got := w.sent.of(team.NotifyObjected); len(got) != 1 || got[0].UserIDs[0] != a.UserID {
		t.Fatalf("objected = %+v", got)
	}
	if err := settler.SettleAt(w.ctx, time.Now().Add(25*time.Hour)); err != nil {
		t.Fatalf("settle: %v", err)
	}
	if n := len(w.scopes.deliveries()); n != 0 {
		t.Fatalf("objection ignored: %d", n)
	}
	if err := settler.SettleAt(w.ctx, due.Add(time.Minute)); err != nil {
		t.Fatalf("settle at due: %v", err)
	}
	got := w.scopes.deliveries()
	if len(got) != 1 {
		t.Fatalf("delivered at due = %+v", got)
	}
	if v := got[0].Votes[1]; v.UserID != b.UserID || v.Decision != module.VoteObjected || v.Comment != "Нет выводов" {
		t.Fatalf("objection not attached: %+v", got[0].Votes)
	}

	d, err = w.teams.Propose(w.ctx, b, d.ID, work("Исправлено"))
	if err != nil {
		t.Fatalf("propose again: %v", err)
	}
	if err := settler.SettleAt(w.ctx, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("settle early: %v", err)
	}
	if err := settler.SettleAt(w.ctx, time.Now().Add(25*time.Hour)); err != nil {
		t.Fatalf("settle silent: %v", err)
	}
	got = w.scopes.deliveries()
	if len(got) != 2 || got[1].Votes[0].Decision != module.VoteSilent {
		t.Fatalf("silence = %+v", got)
	}
}

func TestAISuggestsTeamsOnlyFromWhatStudentsWrote(t *testing.T) {
	w := newWorld(t)
	teacher, a, b, c, d := w.user("teacher"), w.user("aigerim"), w.user("erlan"), w.user("dana"), w.user("timur")
	w.scopes.set(module.TeamScope{ID: "p3", Title: "Проект", Leads: ids(teacher), Members: ids(a, b, c, d), MinSize: 2, MaxSize: 2, SelfFormed: true})
	if _, err := w.teams.SaveCard(w.ctx, a, team.Card{Skills: "Go, SQL", Hours: "вечером"}); err != nil {
		t.Fatalf("card: %v", err)
	}
	if _, err := w.teams.SaveCard(w.ctx, b, team.Card{Interests: "дизайн", Role: "презентация"}); err != nil {
		t.Fatalf("card: %v", err)
	}
	if _, err := w.teams.SaveCard(w.ctx, c, team.Card{Skills: "статистика"}); err != nil {
		t.Fatalf("card: %v", err)
	}
	card, err := w.teams.Card(w.ctx, a)
	if err != nil || card.Skills != "Go, SQL" || card.UpdatedAt == nil {
		t.Fatalf("my card = %+v, %v", card, err)
	}
	w.model.reply = `{"teams":[{"members":["` + a.UserID + `","` + b.UserID + `"],"why":"Код и презентация дополняют друг друга"},{"members":["` + c.UserID + `","` + d.UserID + `"],"why":"нет"},{"members":["` + c.UserID + `"],"why":"один"}]}`
	got, err := w.teams.Suggest(w.ctx, teacher, "study.assignment", "p3")
	if err != nil || len(got.Teams) != 1 || len(got.Teams[0].Members) != 2 || len(got.Without) != 1 || got.Without[0].UserID != d.UserID {
		t.Fatalf("lead suggestion = %+v, %v", got, err)
	}
	asked := w.model.asked()
	if !strings.Contains(asked, "Go, SQL") || strings.Contains(asked, "timur") || strings.Contains(asked, "@kaznu.kz") {
		t.Fatalf("model got:\n%s", asked)
	}

	if _, err := w.teams.Seek(w.ctx, b, "study.assignment", "p3", team.SeekingInput{Seeking: true}); err != nil {
		t.Fatalf("seek: %v", err)
	}
	w.model.reply = `{"teams":[{"members":["` + a.UserID + `","` + b.UserID + `"],"why":"Подходите по времени"}]}`
	got, err = w.teams.Suggest(w.ctx, a, "study.assignment", "p3")
	if err != nil || len(got.Teams) != 1 || got.Teams[0].Why != "Подходите по времени" {
		t.Fatalf("student suggestion = %+v, %v", got, err)
	}
	w.model.reply = "не JSON"
	if _, err := w.teams.Suggest(w.ctx, a, "study.assignment", "p3"); !errors.Is(err, team.ErrAIFailed) {
		t.Fatalf("bad reply: %v", err)
	}
}
