package league_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/league"
)

func (w *world) task(title string) league.Task {
	w.t.Helper()
	due := time.Now().Add(72 * time.Hour)
	t, err := w.leagues.CreateTask(w.ctx, w.teacher, w.view(w.teacher).ID, league.TaskInput{
		Title: title, Partner: "Kaspi Lab", Points: 300, MinSize: 2, MaxSize: 3, SelfFormed: true, DueAt: &due,
	})
	if err != nil {
		w.t.Fatalf("task: %v", err)
	}
	return t
}

func (w *world) scopes() module.TeamScopeSource {
	return league.Module{}.TeamScopeSources(w.kernel)[0]
}

func (w *world) deliver(t league.Task, team, proposal, text string, members ...auth.Principal) {
	w.t.Helper()
	work, err := json.Marshal(league.SolutionWork{Text: text})
	if err != nil {
		w.t.Fatalf("work: %v", err)
	}
	d := module.TeamDelivery{
		TenantID: w.tenant, ScopeID: t.ID, TeamID: team, TeamName: "Команда " + team, ProposalID: proposal,
		ProposedBy: members[0].UserID, Work: work,
	}
	for _, m := range members {
		d.Votes = append(d.Votes, module.TeamVote{UserID: m.UserID, Decision: module.VoteAgreed})
	}
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		return w.scopes().DeliverTeamWork(w.ctx, tx, d)
	})
	if err != nil {
		w.t.Fatalf("deliver: %v", err)
	}
}

func (w *world) prepare(t league.Task, p auth.Principal, work string) error {
	return w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		_, err := w.scopes().PrepareTeamWork(w.ctx, tx, t.ID, p, json.RawMessage(work))
		return err
	})
}

func TestTeamsEarnTaskPointsAfterDefense(t *testing.T) {
	w := newWorld(t)
	task := w.task("Чат-бот для приёмной комиссии")
	if n := len(w.sent.of(league.NotifyTask)); n != 1 {
		t.Fatalf("task notified %d times", n)
	}
	if _, err := w.leagues.CreateTask(w.ctx, w.aigerim, task.LeagueID, league.TaskInput{Title: "x", Points: 1, MinSize: 2, MaxSize: 2}); !errors.Is(err, league.ErrLeagueNotFound) {
		t.Fatalf("student creates task: %v", err)
	}
	var scope module.TeamScope
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		scope, err = w.scopes().TeamScope(w.ctx, tx, task.ID)
		return err
	})
	if err != nil || scope.Kind != league.TaskScopeKind || scope.MaxSize != 3 || len(scope.Members) != 3 || scope.Context != "Основы биологии" {
		t.Fatalf("scope = %+v, %v", scope, err)
	}
	if err := w.prepare(task, w.aigerim, `{"text":"Бот на базе навигатора","links":["https://git.kaznu.kz/bot"]}`); err != nil {
		t.Fatalf("prepare: %v", err)
	}
	var problem module.TeamWorkProblem
	if err := w.prepare(task, w.aigerim, `{"text":"","links":["ftp://x"]}`); !errors.As(err, &problem) || problem.Code != league.CodeInvalidSolution {
		t.Fatalf("bad work: %v", err)
	}

	w.deliver(task, "t1", "p1", "Бот на базе навигатора", w.aigerim, w.erlan)
	w.deliver(task, "t1", "p1", "Бот на базе навигатора", w.aigerim, w.erlan)
	if n := len(w.sent.of(league.NotifySolution)); n != 1 {
		t.Fatalf("solution notified %d times", n)
	}
	details, err := w.leagues.Task(w.ctx, w.teacher, task.ID)
	if err != nil || len(details.Solutions) != 1 || details.Role != league.RoleLead || details.Solutions[0].Work.Text != "Бот на базе навигатора" {
		t.Fatalf("teacher task = %+v, %v", details, err)
	}
	if len(details.People) != 2 || details.Kind != "study.course" {
		t.Fatalf("task people = %+v, kind %q", details.People, details.Kind)
	}
	if mine, _ := w.leagues.Task(w.ctx, w.dana, task.ID); len(mine.Solutions) != 0 {
		t.Fatalf("other student sees solution: %+v", mine.Solutions)
	}
	solution := details.Solutions[0].ID

	if _, err := w.leagues.Review(w.ctx, w.teacher, solution, league.ReviewInput{Points: 301, Passed: []string{w.aigerim.UserID}}); !errors.Is(err, league.ErrInvalidInput) {
		t.Fatalf("too many points: %v", err)
	}
	if _, err := w.leagues.Review(w.ctx, w.teacher, solution, league.ReviewInput{Points: 200, Passed: []string{w.dana.UserID}}); !errors.Is(err, league.ErrInvalidInput) {
		t.Fatalf("stranger passed: %v", err)
	}
	reviewed, err := w.leagues.Review(w.ctx, w.teacher, solution, league.ReviewInput{Points: 200, Passed: []string{w.aigerim.UserID}, Comment: "Ерлан защитит позже"})
	if err != nil || *reviewed.Points != 200 {
		t.Fatalf("review = %+v, %v", reviewed, err)
	}
	me := w.view(w.aigerim)
	if me.Me.Points != 200 || me.Me.Division != 1 || len(me.Entries) != 1 || me.Entries[0].Kind != league.EntryTask || me.Entries[0].Reason != "Ерлан защитит позже" {
		t.Fatalf("aigerim = %+v %+v", me.Me, me.Entries)
	}
	if he := w.view(w.erlan); he.Me.Points != 0 {
		t.Fatalf("erlan without defense = %+v", he.Me)
	}
	if _, err := w.leagues.Review(w.ctx, w.teacher, solution, league.ReviewInput{Points: 250, Passed: []string{w.aigerim.UserID, w.erlan.UserID}}); err != nil {
		t.Fatalf("second review: %v", err)
	}
	if he := w.view(w.erlan); he.Me.Points != 250 || len(he.Entries) != 1 {
		t.Fatalf("erlan after defense = %+v", he.Me)
	}
	if n := len(w.sent.of(league.NotifyReviewed)); n != 2 {
		t.Fatalf("reviewed notified %d times", n)
	}

	task.ClosedAt = nil
	if _, err := w.leagues.UpdateTask(w.ctx, w.teacher, task.ID, league.TaskInput{Title: task.Title, Points: 300, MinSize: 2, MaxSize: 3, Closed: true}); err != nil {
		t.Fatalf("close task: %v", err)
	}
	if err := w.prepare(task, w.aigerim, `{"text":"Ещё версия"}`); !errors.As(err, &problem) || problem.Code != league.CodeTaskClosed {
		t.Fatalf("closed task: %v", err)
	}
	if v := w.view(w.dana); len(v.Tasks) != 1 || v.Tasks[0].ClosedAt == nil {
		t.Fatalf("tasks in view = %+v", v.Tasks)
	}
}

func TestSimilarSolutionsAreShownToTeacher(t *testing.T) {
	w := newWorld(t)
	first := w.task("Карта кампуса")
	second := w.task("Навигация по кампусу")
	w.deliver(first, "t1", "p1", "Интерактивная карта корпусов с поиском аудиторий", w.aigerim, w.erlan)
	w.deliver(second, "t2", "p2", "Интерактивная карта корпусов с поиском аудиторий", w.dana, w.erlan)
	w.deliver(second, "t3", "p3", "Расписание", w.aigerim, w.dana)
	details, err := w.leagues.Task(w.ctx, w.teacher, second.ID)
	if err != nil || len(details.Solutions) != 2 {
		t.Fatalf("details = %+v, %v", details, err)
	}
	flagged := details.Solutions[0]
	if flagged.Similar == nil || flagged.Similar.TaskTitle != "Карта кампуса" || *flagged.Similarity < 0.99 || details.Solutions[1].Similar != nil {
		t.Fatalf("similar = %+v / %+v", flagged, details.Solutions[1])
	}
	if _, err := w.leagues.Explain(w.ctx, w.aigerim, flagged.ID, league.ExplanationInput{Text: "Не моя команда"}); !errors.Is(err, league.ErrSolutionNotFound) {
		t.Fatalf("outsider explains: %v", err)
	}
	explained, err := w.leagues.Explain(w.ctx, w.dana, flagged.ID, league.ExplanationInput{Text: "Мы добавили маршруты для колясок"})
	if err != nil || explained.Explanation == "" {
		t.Fatalf("explain = %+v, %v", explained, err)
	}
	mine, err := w.leagues.Task(w.ctx, w.dana, second.ID)
	if err != nil || len(mine.Solutions) != 2 || mine.Solutions[0].Similar != nil {
		t.Fatalf("student sees similar project: %+v, %v", mine.Solutions, err)
	}
}

func TestUniversityAwardsAndDivisionQuota(t *testing.T) {
	w := newWorld(t)
	one := 1
	if _, err := w.leagues.CreateAward(w.ctx, w.teacher, league.AwardInput{Title: "Часы GPU", GPUHours: 20}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher creates award: %v", err)
	}
	award, err := w.leagues.CreateAward(w.ctx, w.admin, league.AwardInput{Title: "Стажировка в Kaspi Lab", Quantity: &one})
	if err != nil {
		t.Fatalf("award: %v", err)
	}
	id := w.view(w.teacher).ID
	grant, err := w.leagues.Grant(w.ctx, w.teacher, id, league.GrantInput{AwardID: award.ID, UserID: w.aigerim.UserID, Reason: "Самый большой рост"})
	if err != nil || grant.Title != "Стажировка в Kaspi Lab" {
		t.Fatalf("grant = %+v, %v", grant, err)
	}
	if _, err := w.leagues.Grant(w.ctx, w.teacher, id, league.GrantInput{AwardID: award.ID, UserID: w.erlan.UserID, Reason: "Тоже"}); !errors.Is(err, league.ErrAwardGone) {
		t.Fatalf("second grant: %v", err)
	}
	if _, err := w.leagues.Grant(w.ctx, w.aigerim, id, league.GrantInput{AwardID: award.ID, UserID: w.aigerim.UserID, Reason: "Сама"}); !errors.Is(err, league.ErrLeagueNotFound) {
		t.Fatalf("student grants: %v", err)
	}
	if n := len(w.sent.of(league.NotifyAwarded)); n != 1 {
		t.Fatalf("awarded notified %d times", n)
	}
	if v := w.view(w.aigerim); len(v.Grants) != 1 || v.Grants[0].Reason != "Самый большой рост" {
		t.Fatalf("student grants = %+v", v.Grants)
	}
	if v := w.view(w.erlan); len(v.Grants) != 0 {
		t.Fatalf("other grants leaked = %+v", v.Grants)
	}
	if err := w.leagues.Revoke(w.ctx, w.teacher, grant.ID); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := w.leagues.Grant(w.ctx, w.teacher, id, league.GrantInput{AwardID: award.ID, UserID: w.erlan.UserID, Reason: "Защита лучше всех"}); err != nil {
		t.Fatalf("grant after revoke: %v", err)
	}
	list, err := w.leagues.Awards(w.ctx, w.aigerim)
	if err != nil || len(list) != 1 || list[0].Granted != 1 {
		t.Fatalf("awards = %+v, %v", list, err)
	}

	quota := league.Divisions{{From: 0}, {From: 150, GPUHours: 10}}
	if _, err := w.leagues.SetRules(w.ctx, w.teacher, id, league.RulesInput{Divisions: quota}); !errors.Is(err, league.ErrInvalidRules) {
		t.Fatalf("quota over university limit: %v", err)
	}
	w.settings.set(league.DivisionGPUHours, 10)
	w.rules(league.RulesInput{Divisions: quota, Weights: weights(200, 0)})
	w.course.grade("project", w.dana, 10, 10)
	if v := w.view(w.dana); v.Me.Division != 1 || v.Me.GPUHours != 10 {
		t.Fatalf("quota = %+v", v.Me)
	}

	gpuAward, err := w.leagues.CreateAward(w.ctx, w.admin, league.AwardInput{Title: "Часы GPU", GPUHours: 20})
	if err != nil {
		t.Fatalf("gpu award: %v", err)
	}
	if _, err := w.leagues.Grant(w.ctx, w.teacher, id, league.GrantInput{AwardID: gpuAward.ID, UserID: w.dana.UserID, Reason: "Лучший проект"}); err != nil {
		t.Fatalf("gpu grant: %v", err)
	}
	var grants []module.GPUGrant
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		grants, err = league.GPUGrants{}.GPUGrants(w.ctx, tx, w.dana.UserID, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
		return err
	})
	want := []module.GPUGrant{
		{Source: league.GPUSourceDivision, Title: "2", Hours: 10},
		{Source: league.GPUSourceAward, Title: "Часы GPU", Hours: 20},
	}
	if err != nil || len(grants) != 2 || grants[0].Source != want[0].Source || grants[0].Hours != 10 || grants[1] != want[1] {
		t.Fatalf("gpu grants = %+v, %v", grants, err)
	}

	rec := w.call(w.teacher, http.MethodDelete, "/v1/league/grants/"+grant.ID, "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("revoke twice = %d %s", rec.Code, rec.Body)
	}
}
