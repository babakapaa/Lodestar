package sandbox_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/sandbox"
)

func TestAdministratorSetsClassesAndTerms(t *testing.T) {
	w := newWorld(t)
	if _, err := w.quotas.CreateClass(w.ctx, w.teacher, sandbox.GPUClassInput{Name: "T4", MemoryGB: 16, Rate: 1}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher adds class: %v", err)
	}
	if _, err := w.quotas.CreateClass(w.ctx, w.admin, sandbox.GPUClassInput{Name: "T4", MemoryGB: 16, Rate: 0}); !errors.Is(err, sandbox.ErrInvalidClass) {
		t.Fatalf("free class: %v", err)
	}
	t4, err := w.quotas.CreateClass(w.ctx, w.admin, sandbox.GPUClassInput{Name: "T4", MemoryGB: 16, Rate: 1})
	if err != nil {
		t.Fatalf("t4: %v", err)
	}
	if _, err := w.quotas.CreateClass(w.ctx, w.admin, sandbox.GPUClassInput{Name: "A100", MemoryGB: 80, Rate: 4.555}); err != nil {
		t.Fatalf("a100: %v", err)
	}
	if _, err := w.quotas.CreateClass(w.ctx, w.admin, sandbox.GPUClassInput{Name: "T4", MemoryGB: 16, Rate: 1}); !errors.Is(err, sandbox.ErrClassTaken) {
		t.Fatalf("same name: %v", err)
	}
	if _, err := w.quotas.UpdateClass(w.ctx, w.admin, t4.ID, sandbox.GPUClassInput{Name: "T4", MemoryGB: 16, Rate: 1, Archived: true}); err != nil {
		t.Fatalf("archive: %v", err)
	}
	list, err := w.quotas.Classes(w.ctx, w.aigerim)
	if err != nil || len(list) != 2 || list[0].Name != "A100" || list[0].Rate != 4.56 || !list[1].Archived {
		t.Fatalf("classes = %+v, %v", list, err)
	}

	autumn := w.term()
	if _, err := w.quotas.CreateTerm(w.ctx, w.admin, sandbox.TermInput{Name: "Лето", StartsOn: autumn.EndsOn, EndsOn: "2099-01-01"}); !errors.Is(err, sandbox.ErrTermOverlap) {
		t.Fatalf("overlap: %v", err)
	}
	if _, err := w.quotas.CreateTerm(w.ctx, w.admin, sandbox.TermInput{Name: "Наоборот", StartsOn: "2099-02-01", EndsOn: "2099-01-01"}); !errors.Is(err, sandbox.ErrInvalidTerm) {
		t.Fatalf("backwards: %v", err)
	}
	rec := w.call(w.admin, http.MethodPost, "/v1/sandbox/terms", `{"name":"Весна","startsOn":"2099-01-10","endsOn":"2099-05-31"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("spring = %d %s", rec.Code, rec.Body)
	}
	if terms, err := w.quotas.Terms(w.ctx, w.admin); err != nil || len(terms) != 2 || terms[0].Name != "Весна" {
		t.Fatalf("terms = %+v, %v", terms, err)
	}
}

func TestQuotaAddsBaseLeagueAndTeam(t *testing.T) {
	w := newWorld(t)
	if q := w.quota(w.aigerim); q.Term != nil || q.Hours != 0 {
		t.Fatalf("quota without term = %+v", q)
	}
	if _, err := w.quotas.Transfer(w.ctx, w.aigerim, sandbox.TransferInput{TeamID: "team-1", Hours: 1}); !errors.Is(err, sandbox.ErrNoTerm) {
		t.Fatalf("transfer without term: %v", err)
	}
	w.term()
	w.settings.set(sandbox.BaseHours, 30)
	q := w.quota(w.aigerim)
	if q.Hours != 40 || len(q.Lines) != 2 || q.Lines[0].Source != sandbox.SourceBase || q.Lines[1].Hours != 10 || len(q.Teams) != 1 {
		t.Fatalf("aigerim quota = %+v", q)
	}
	if q := w.quota(w.teacher); q.Hours != 30 || len(q.Teams) != 0 {
		t.Fatalf("teacher quota = %+v", q)
	}

	if _, err := w.quotas.Transfer(w.ctx, w.aigerim, sandbox.TransferInput{TeamID: "team-2", Hours: 1}); !errors.Is(err, sandbox.ErrNotInTeam) {
		t.Fatalf("foreign team: %v", err)
	}
	if _, err := w.quotas.Transfer(w.ctx, w.aigerim, sandbox.TransferInput{TeamID: "team-1", Hours: 41}); !errors.Is(err, sandbox.ErrNotEnoughHours) {
		t.Fatalf("too much: %v", err)
	}
	q, err := w.quotas.Transfer(w.ctx, w.aigerim, sandbox.TransferInput{TeamID: "team-1", Hours: 12.5})
	if err != nil || q.Hours != 27.5 || q.Teams[0].Hours != 12.5 || q.Teams[0].Lines[0].Title != "Айгерим" {
		t.Fatalf("after transfer = %+v, %v", q, err)
	}
	if q := w.quota(w.erlan); len(q.Teams) != 1 || q.Teams[0].Hours != 12.5 {
		t.Fatalf("teammate sees pool = %+v", q)
	}
	rec := w.call(w.erlan, http.MethodPost, "/v1/sandbox/transfers", `{"teamId":"team-1","hours":-5}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("negative transfer = %d %s", rec.Code, rec.Body)
	}
}

func TestTopupGoesToTeacher(t *testing.T) {
	w := newWorld(t)
	w.term()
	in := sandbox.TopupInput{WorkKind: "assignment", WorkID: "a1", Hours: 8, Reason: "Не успеваю обучить модель"}
	if _, err := w.quotas.RequestTopup(w.ctx, w.erlan, in); !errors.Is(err, sandbox.ErrWorkNotFound) {
		t.Fatalf("stranger asks: %v", err)
	}
	if _, err := w.quotas.RequestTopup(w.ctx, w.teacher, in); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher asks: %v", err)
	}
	topup, err := w.quotas.RequestTopup(w.ctx, w.aigerim, in)
	if err != nil || topup.Status != sandbox.TopupPending || topup.Title != "Нейросеть для МРТ" {
		t.Fatalf("request = %+v, %v", topup, err)
	}
	if _, err := w.quotas.RequestTopup(w.ctx, w.aigerim, in); !errors.Is(err, sandbox.ErrTopupPending) {
		t.Fatalf("second request: %v", err)
	}
	asked := w.sent.of(sandbox.NotifyTopupRequested)
	if len(asked) != 1 || asked[0].UserIDs[0] != w.teacher.UserID || asked[0].Params["hours"] != "8" {
		t.Fatalf("teacher told = %+v", asked)
	}

	if list, err := w.quotas.Topups(w.ctx, w.teacher); err != nil || len(list) != 1 || list[0].Mine || list[0].DisplayName != "Айгерим" {
		t.Fatalf("teacher list = %+v, %v", list, err)
	}
	if list, err := w.quotas.Topups(w.ctx, w.erlan); err != nil || len(list) != 0 {
		t.Fatalf("stranger list = %+v, %v", list, err)
	}
	if _, err := w.quotas.DecideTopup(w.ctx, w.aigerim, topup.ID, sandbox.TopupDecision{Approve: true}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student approves: %v", err)
	}
	five := 5.0
	decided, err := w.quotas.DecideTopup(w.ctx, w.teacher, topup.ID, sandbox.TopupDecision{Approve: true, Hours: &five})
	if err != nil || decided.Status != sandbox.TopupApproved || *decided.Granted != 5 || decided.DecidedAt == nil {
		t.Fatalf("approve = %+v, %v", decided, err)
	}
	if _, err := w.quotas.DecideTopup(w.ctx, w.teacher, topup.ID, sandbox.TopupDecision{}); !errors.Is(err, sandbox.ErrTopupNotFound) {
		t.Fatalf("decide twice: %v", err)
	}
	q := w.quota(w.aigerim)
	if q.Hours != 35 || q.Lines[2].Source != sandbox.SourceTopup || q.Lines[2].Title != "Нейросеть для МРТ" {
		t.Fatalf("quota after top-up = %+v", q)
	}
	if n := len(w.sent.of(sandbox.NotifyTopupApproved)); n != 1 {
		t.Fatalf("student told %d times", n)
	}

	again, err := w.quotas.RequestTopup(w.ctx, w.aigerim, in)
	if err != nil {
		t.Fatalf("request again: %v", err)
	}
	rec := w.call(w.teacher, http.MethodPost, "/v1/sandbox/topups/"+again.ID+"/decision", `{"approve":false}`)
	if rec.Code != http.StatusOK || len(w.sent.of(sandbox.NotifyTopupDeclined)) != 1 {
		t.Fatalf("decline = %d %s", rec.Code, rec.Body)
	}
	if q := w.quota(w.aigerim); q.Hours != 35 {
		t.Fatalf("declined adds hours = %+v", q)
	}
}
