package league_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/league"
)

func TestPointsFollowGradesByVisibleFormula(t *testing.T) {
	w := newWorld(t)
	w.course.grade("project", w.aigerim, 8, 10)
	w.course.grade("project", w.erlan, 10, 10)
	w.course.grade("quiz", w.dana, 1, 1)

	v := w.view(w.aigerim)
	if v.Role != league.RoleStudent || v.Season.Number != 1 || v.Me == nil || v.Me.Points != 80 || v.Me.Division != 0 || *v.Me.ToNext != 70 {
		t.Fatalf("student view = %+v, me %+v", v, v.Me)
	}
	if len(v.Formula) != 2 || v.Formula[0].Weight != 100 || v.Formula[1].Weight != 0 || v.Formula[1].Project {
		t.Fatalf("formula = %+v", v.Formula)
	}
	if len(v.Table) != 3 || v.Table[0].Name != "Ерлан" || v.Table[0].Place != 1 || !v.Table[1].Me || v.Table[2].Points != 0 {
		t.Fatalf("table = %+v", v.Table)
	}
	if len(v.Entries) != 1 || *v.Entries[0].Score != 8 || *v.Entries[0].Weight != 100 || v.Entries[0].Title != "Проект" {
		t.Fatalf("entries = %+v", v.Entries)
	}
	if v.Players != nil || v.Pending != nil {
		t.Fatalf("student sees teacher data: %+v", v)
	}

	v = w.rules(league.RulesInput{Divisions: three(), Weights: weights(200, 50)})
	if v.Role != league.RoleLead || len(v.Players) != 3 || v.Players[0].Division != 1 || v.Players[0].Points != 200 {
		t.Fatalf("teacher view = %+v", v.Players)
	}
	promoted := w.sent.of(league.NotifyPromoted)
	if len(promoted) != 2 || promoted[0].Params["division"] != "2" {
		t.Fatalf("promoted = %+v", promoted)
	}
	me := w.view(w.aigerim)
	if me.Me.Points != 160 || me.Me.Division != 1 || *me.Me.ToNext != 140 || len(me.Table) != 2 || len(me.Transitions) != 1 {
		t.Fatalf("after rules = %+v, %+v", me.Me, me.Table)
	}
	if d := w.view(w.dana); d.Me.Points != 50 || d.Me.Division != 0 || len(d.Table) != 1 {
		t.Fatalf("dana = %+v", d.Me)
	}

	w.course.grade("project", w.aigerim, 10, 10)
	if me := w.view(w.aigerim); me.Me.Points != 200 || *me.Entries[0].Score != 10 {
		t.Fatalf("regrade = %+v", me.Me)
	}

	rec := w.call(w.erlan, http.MethodGet, "/v1/league/arenas/study.course/"+course, "")
	var body map[string]any
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &body) != nil || body["role"] != "student" {
		t.Fatalf("http = %d %s", rec.Code, rec.Body)
	}
}

func TestNamesCanBeHiddenOrReplaced(t *testing.T) {
	w := newWorld(t)
	w.course.grade("project", w.aigerim, 5, 10)
	id := w.view(w.aigerim).ID
	if _, err := w.leagues.SetDisplay(w.ctx, w.aigerim, id, league.DisplayInput{Display: league.DisplayHidden}); err != nil {
		t.Fatalf("hide: %v", err)
	}
	if _, err := w.leagues.SetDisplay(w.ctx, w.erlan, id, league.DisplayInput{Display: league.DisplayNickname, Nickname: " Ёж "}); err != nil {
		t.Fatalf("nickname: %v", err)
	}
	if _, err := w.leagues.SetDisplay(w.ctx, w.dana, id, league.DisplayInput{Display: league.DisplayNickname}); !errors.Is(err, league.ErrInvalidInput) {
		t.Fatalf("empty nickname: %v", err)
	}
	v := w.view(w.dana)
	names := map[string]league.Display{}
	for _, row := range v.Table {
		names[row.Name] = row.Display
	}
	if names[""] != league.DisplayHidden || names["Ёж"] != league.DisplayNickname || names["Дана"] != league.DisplayName || len(names) != 3 {
		t.Fatalf("table = %+v", v.Table)
	}
	if v.Table[0].Points != 50 || v.Table[0].Name != "" {
		t.Fatalf("hidden leader = %+v", v.Table[0])
	}
	teacher := w.view(w.teacher)
	for _, p := range teacher.Players {
		if p.Person.DisplayName == "" {
			t.Fatalf("teacher sees hidden names: %+v", teacher.Players)
		}
	}
}

func TestQueuedMovesAndSeasons(t *testing.T) {
	w := newWorld(t)
	w.rules(league.RulesInput{Divisions: three(), Confirm: true, Weights: weights(400, 0)})
	w.course.grade("project", w.erlan, 10, 10)
	w.course.grade("project", w.aigerim, 8, 10)
	v := w.view(w.teacher)
	if len(v.Pending) != 2 || len(w.sent.of(league.NotifyQueued)) != 1 {
		t.Fatalf("pending = %+v", v.Pending)
	}
	if len(w.sent.of(league.NotifyPromoted)) != 0 {
		t.Fatal("promoted without confirmation")
	}
	for _, p := range v.Pending {
		approve := p.UserID == w.erlan.UserID
		if _, err := w.leagues.Decide(w.ctx, w.teacher, p.ID, league.DecisionInput{Approve: approve}); err != nil {
			t.Fatalf("decide: %v", err)
		}
		if _, err := w.leagues.Decide(w.ctx, w.teacher, p.ID, league.DecisionInput{Approve: approve}); !errors.Is(err, league.ErrTransitionNotFound) {
			t.Fatalf("decide twice: %v", err)
		}
	}
	v = w.view(w.teacher)
	if len(v.Pending) != 0 || v.Players[0].Person.UserID != w.erlan.UserID || v.Players[0].Division != 2 {
		t.Fatalf("after decisions = %+v / %+v", v.Pending, v.Players)
	}
	if me := w.view(w.aigerim); me.Me.Division != 0 || me.Transitions[0].Status != league.TransitionRejected {
		t.Fatalf("rejected = %+v", me)
	}

	v, err := w.leagues.CloseSeason(w.ctx, w.teacher, v.ID)
	if err != nil || v.Season.Number != 2 || len(v.Pending) != 0 {
		t.Fatalf("close 1 = %+v, %v", v.Season, err)
	}
	w.course.grade("project", w.dana, 5, 10)
	me := w.view(w.erlan)
	if me.Me.Points != 0 || me.Me.Division != 2 || me.Me.Growth != -400 || len(me.Entries) != 0 {
		t.Fatalf("season 2 = %+v", me.Me)
	}
	dana := w.view(w.dana)
	if dana.Me.Points != 200 || dana.Me.Growth != 200 || !dana.Table[0].TopGrowth {
		t.Fatalf("dana season 2 = %+v %+v", dana.Me, dana.Table)
	}
	w.course.grade("project", w.erlan, 1, 10)
	if me := w.view(w.erlan); me.Me.Points != 0 {
		t.Fatalf("closed season entry moved: %+v", me.Me)
	}

	v, err = w.leagues.CloseSeason(w.ctx, w.teacher, v.ID)
	if err != nil || v.Season.Number != 3 || len(v.Pending) != 2 {
		t.Fatalf("close 2 = %+v, %v", v.Pending, err)
	}
	down := v.Pending[1]
	if down.UserID != w.erlan.UserID || down.From != 2 || down.To != 0 || v.Pending[0].UserID != w.dana.UserID {
		t.Fatalf("pending after close = %+v", v.Pending)
	}
	if _, err := w.leagues.Decide(w.ctx, w.teacher, down.ID, league.DecisionInput{Approve: true}); err != nil {
		t.Fatalf("approve down: %v", err)
	}
	if down := w.sent.of(league.NotifyRelegated); len(down) != 1 || down[0].UserIDs[0] != w.erlan.UserID {
		t.Fatalf("relegated = %+v", down)
	}
	if me := w.view(w.erlan); me.Me.Division != 0 {
		t.Fatalf("after down = %+v", me.Me)
	}
}

func TestPointsAndMovesCanBeDisputed(t *testing.T) {
	w := newWorld(t)
	w.course.grade("project", w.aigerim, 7, 10)
	w.course.grade("project", w.erlan, 9, 10)
	mine := w.view(w.aigerim)
	theirs := w.view(w.erlan)
	entry := mine.Entries[0].ID
	if _, err := w.leagues.Dispute(w.ctx, w.erlan, mine.ID, league.DisputeInput{EntryID: &entry, Text: "Не моё"}); !errors.Is(err, league.ErrDisputeTarget) {
		t.Fatalf("foreign entry: %v", err)
	}
	d, err := w.leagues.Dispute(w.ctx, w.aigerim, mine.ID, league.DisputeInput{EntryID: &entry, Text: "Не учли доработку"})
	if err != nil || d.Status != league.DisputeOpen {
		t.Fatalf("dispute = %+v, %v", d, err)
	}
	if _, err := w.leagues.Dispute(w.ctx, w.aigerim, mine.ID, league.DisputeInput{EntryID: &entry, Text: "Ещё раз"}); !errors.Is(err, league.ErrDisputeOpen) {
		t.Fatalf("second dispute: %v", err)
	}
	if n := len(w.sent.of(league.NotifyDisputed)); n != 1 {
		t.Fatalf("disputed notified %d times", n)
	}
	teacher := w.view(w.teacher)
	if len(teacher.Open) != 1 || teacher.Open[0].Entry == nil || teacher.Open[0].Person.UserID != w.aigerim.UserID {
		t.Fatalf("open = %+v", teacher.Open)
	}
	if _, err := w.leagues.Answer(w.ctx, w.aigerim, d.ID, league.AnswerInput{Accept: true, Answer: "Сама себе"}); !errors.Is(err, league.ErrDisputeNotFound) {
		t.Fatalf("student answers: %v", err)
	}
	answered, err := w.leagues.Answer(w.ctx, w.teacher, d.ID, league.AnswerInput{Accept: true, Answer: "Доработка учтена", Points: 80})
	if err != nil || answered.Status != league.DisputeAccepted {
		t.Fatalf("answer = %+v, %v", answered, err)
	}
	mine = w.view(w.aigerim)
	if mine.Me.Points != 150 || mine.Me.Division != 1 || len(mine.Entries) != 2 || mine.Entries[1].Reason != "Доработка учтена" {
		t.Fatalf("after answer = %+v %+v", mine.Me, mine.Entries)
	}
	if n := len(w.sent.of(league.NotifyAnswered)); n != 1 {
		t.Fatalf("answered notified %d times", n)
	}

	move := mine.Transitions[0].ID
	d, err = w.leagues.Dispute(w.ctx, w.aigerim, mine.ID, league.DisputeInput{TransitionID: &move, Text: "Хочу заработать сама"})
	if err != nil {
		t.Fatalf("dispute move: %v", err)
	}
	if _, err := w.leagues.Answer(w.ctx, w.teacher, d.ID, league.AnswerInput{Accept: true, Answer: "Вернула", Points: 5}); !errors.Is(err, league.ErrInvalidInput) {
		t.Fatalf("points on move: %v", err)
	}
	if _, err := w.leagues.Answer(w.ctx, w.teacher, d.ID, league.AnswerInput{Accept: true, Answer: "Вернула"}); err != nil {
		t.Fatalf("revert: %v", err)
	}
	if mine = w.view(w.aigerim); mine.Me.Division != 0 || mine.Transitions[0].Status != league.TransitionReverted {
		t.Fatalf("after revert = %+v %+v", mine.Me, mine.Transitions)
	}

	adjusted, err := w.leagues.Adjust(w.ctx, w.teacher, theirs.ID, league.AdjustmentInput{UserID: w.erlan.UserID, Points: -10, Reason: "Повтор прошлогоднего проекта"})
	if err != nil || adjusted.Standing.Points != 80 || len(adjusted.Entries) != 2 {
		t.Fatalf("adjust = %+v, %v", adjusted, err)
	}
	if _, err := w.leagues.Adjust(w.ctx, w.teacher, theirs.ID, league.AdjustmentInput{UserID: w.erlan.UserID, Points: 10}); !errors.Is(err, league.ErrInvalidInput) {
		t.Fatalf("adjust without reason: %v", err)
	}
}

func TestLeagueIsClosedToStrangersAndStudentsCannotRule(t *testing.T) {
	w := newWorld(t)
	stranger := w.user("stranger", "Чужой")
	if _, err := w.leagues.View(w.ctx, stranger, "study.course", course); !errors.Is(err, league.ErrLeagueNotFound) {
		t.Fatalf("stranger: %v", err)
	}
	if _, err := w.leagues.View(w.ctx, w.aigerim, "study.course", "other"); !errors.Is(err, league.ErrLeagueNotFound) {
		t.Fatalf("other course: %v", err)
	}
	id := w.view(w.aigerim).ID
	if _, err := w.leagues.SetRules(w.ctx, w.aigerim, id, league.RulesInput{Divisions: three()}); !errors.Is(err, league.ErrLeagueNotFound) {
		t.Fatalf("student rules: %v", err)
	}
	if _, err := w.leagues.CloseSeason(w.ctx, w.aigerim, id); !errors.Is(err, league.ErrLeagueNotFound) {
		t.Fatalf("student closes season: %v", err)
	}
	bad := []league.RulesInput{
		{Divisions: league.Divisions{{From: 10}}},
		{Divisions: league.Divisions{{From: 0}, {From: 0}}},
		{Divisions: three(), Weights: []league.WeightInput{{Kind: "study.assignment", ID: "missing", Weight: 10}}},
		{Divisions: three(), Weights: []league.WeightInput{{Kind: "study.assignment", ID: "quiz", Weight: -1}}},
	}
	for _, in := range bad {
		if _, err := w.leagues.SetRules(w.ctx, w.teacher, id, in); !errors.Is(err, league.ErrInvalidRules) {
			t.Fatalf("rules %+v: %v", in, err)
		}
	}
	if v := w.rules(league.RulesInput{Divisions: league.Divisions{{From: 0}, {From: 50}}}); len(v.Divisions) != 2 {
		t.Fatalf("two divisions = %+v", v.Divisions)
	}
}

func TestWorkerSyncsActiveCourses(t *testing.T) {
	w := newWorld(t)
	w.course.grade("project", w.aigerim, 10, 10)
	w.course.grade("project", w.erlan, 10, 10)
	if err := w.worker.SyncTenant(w.ctx, w.tenant); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if n := len(w.sent.of(league.NotifyPromoted)); n != 0 {
		t.Fatalf("promoted %d below boundary", n)
	}
	w.course.grade("quiz", w.aigerim, 1, 1)
	w.rules(league.RulesInput{Divisions: three(), Weights: weights(100, 60)})
	if err := w.worker.SyncTenant(w.ctx, w.tenant); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if up := w.sent.of(league.NotifyPromoted); len(up) != 1 || up[0].UserIDs[0] != w.aigerim.UserID {
		t.Fatalf("promoted = %+v", up)
	}
}
