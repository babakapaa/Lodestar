package sandbox_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/sandbox"
)

func busy(processes ...string) *sandbox.RunUsage {
	return &sandbox.RunUsage{CPU: 80, NetBytes: 1 << 20, Processes: processes}
}

func quiet() *sandbox.RunUsage {
	return &sandbox.RunUsage{CPU: 0.4, NetBytes: 1<<20 + 512}
}

func TestIdleSessionIsWarnedThenStopped(t *testing.T) {
	w := newWorld(t)
	w.term()
	f := w.fleet(1)
	run, err := w.quotas.StartRun(w.ctx, w.aigerim, sandbox.RunInput{Kind: sandbox.RunJupyter, Title: "Ноутбук", TemplateID: f.template.ID, ClassID: f.t4.ID})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	w.sync(f.node.Token)
	w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseRunning, URL: "https://gpu/s", Usage: busy("python")})
	w.exec(`UPDATE sandbox.runs SET active_at = now() - interval '21 minutes' WHERE id = $1`, run.ID)
	w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseRunning, Usage: quiet()})
	w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseRunning, Usage: quiet()})
	warned := w.sent.of(sandbox.NotifyIdle)
	if len(warned) != 1 || warned[0].Params["minutes"] != 9 || warned[0].UserIDs[0] != w.aigerim.UserID {
		t.Fatalf("warned = %+v", warned)
	}
	w.exec(`UPDATE sandbox.runs SET active_at = now() - interval '31 minutes' WHERE id = $1`, run.ID)
	w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseRunning, Usage: quiet()})
	if plan := w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseRunning}); len(plan.Runs) != 1 || !plan.Runs[0].Stop {
		t.Fatalf("idle session kept = %+v", plan)
	}
	w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseExited, ExitCode: 137})
	ended, err := w.quotas.GetRun(w.ctx, w.aigerim, run.ID)
	if err != nil || ended.Status != sandbox.RunStopped || ended.Reason != sandbox.ReasonIdle {
		t.Fatalf("ended = %+v, %v", ended, err)
	}
	if n := len(w.sent.of(sandbox.NotifyRunFailed)); n != 0 {
		t.Fatalf("idle stop notified %d times", n)
	}
}

func TestWatchdogHoldsAndAdminDecides(t *testing.T) {
	w := newWorld(t)
	w.term()
	f := w.fleet(2)
	job := sandbox.RunInput{Kind: sandbox.RunJob, Title: "Обучение", Command: "python train.py", TemplateID: f.template.ID, ClassID: f.t4.ID, MaxHours: 4}
	first, err := w.quotas.StartRun(w.ctx, w.aigerim, job)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	second, err := w.quotas.StartRun(w.ctx, w.aigerim, job)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	w.sync(f.node.Token)
	w.sync(f.node.Token,
		sandbox.RunReport{RunID: first.ID, Phase: sandbox.PhaseRunning, Usage: busy("python")},
		sandbox.RunReport{RunID: second.ID, Phase: sandbox.PhaseRunning, Usage: busy("python")})
	w.exec(`UPDATE sandbox.runs SET started_at = now() - interval '3 hours'`)
	plan := w.sync(f.node.Token,
		sandbox.RunReport{RunID: first.ID, Phase: sandbox.PhaseRunning, Usage: busy("python", "/opt/XMRig")},
		sandbox.RunReport{RunID: second.ID, Phase: sandbox.PhaseRunning, Usage: busy("sh", "xmrig")})
	for _, spec := range plan.Runs {
		if !spec.Hold || spec.Stop {
			t.Fatalf("plan = %+v", plan)
		}
	}
	held, err := w.quotas.GetRun(w.ctx, w.aigerim, first.ID)
	if err != nil || held.Status != sandbox.RunHeld || held.Reason != sandbox.ReasonWatchdog || held.HoldNote != "xmrig" || held.HeldAt == nil {
		t.Fatalf("held = %+v, %v", held, err)
	}
	if _, err := w.quotas.HeldRuns(w.ctx, w.teacher); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher sees held: %v", err)
	}
	if _, err := w.quotas.DecideHold(w.ctx, w.teacher, first.ID, sandbox.HoldDecision{Release: true}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher decides: %v", err)
	}
	list, err := w.quotas.HeldRuns(w.ctx, w.admin)
	if err != nil || len(list) != 2 || list[0].URL != "" {
		t.Fatalf("held list = %+v, %v", list, err)
	}
	if rec := w.call(w.admin, http.MethodGet, "/v1/sandbox/runs/held", ""); rec.Code != http.StatusOK || strings.Count(rec.Body.String(), `"holdNote":"xmrig"`) != 2 {
		t.Fatalf("held over http = %d %s", rec.Code, rec.Body)
	}
	if q := w.quota(w.aigerim); q.Reserved != 16 {
		t.Fatalf("held runs keep reserve = %+v", q)
	}
	w.exec(`UPDATE sandbox.runs SET held_at = now() - interval '2 hours'`)
	if rec := w.call(w.admin, http.MethodPost, "/v1/sandbox/runs/"+first.ID+"/hold-decision", `{"release":true}`); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"status":"running"`) {
		t.Fatalf("release = %d %s", rec.Code, rec.Body)
	}
	if _, err := w.quotas.DecideHold(w.ctx, w.admin, second.ID, sandbox.HoldDecision{}); err != nil {
		t.Fatalf("stop held: %v", err)
	}
	if _, err := w.quotas.DecideHold(w.ctx, w.admin, second.ID, sandbox.HoldDecision{Release: true}); !errors.Is(err, sandbox.ErrRunNotHeld) {
		t.Fatalf("decide twice: %v", err)
	}
	plan = w.sync(f.node.Token,
		sandbox.RunReport{RunID: first.ID, Phase: sandbox.PhaseRunning, Usage: busy("python")},
		sandbox.RunReport{RunID: second.ID, Phase: sandbox.PhaseRunning})
	for _, spec := range plan.Runs {
		if spec.Hold || spec.Stop != (spec.RunID == second.ID) {
			t.Fatalf("plan after decisions = %+v", plan)
		}
	}
	w.sync(f.node.Token,
		sandbox.RunReport{RunID: first.ID, Phase: sandbox.PhaseExited},
		sandbox.RunReport{RunID: second.ID, Phase: sandbox.PhaseExited, ExitCode: 137})
	done, _ := w.quotas.GetRun(w.ctx, w.aigerim, first.ID)
	stopped, _ := w.quotas.GetRun(w.ctx, w.aigerim, second.ID)
	if done.Status != sandbox.RunFinished || *done.Charged != 2 || stopped.Status != sandbox.RunStopped || *stopped.Charged != 2 {
		t.Fatalf("held time charged: done = %+v, stopped = %+v", done, stopped)
	}
	if got := w.sent.of(sandbox.NotifyRunFinished); len(got) != 1 || got[0].Params["hours"] != "2" {
		t.Fatalf("finished = %+v", got)
	}
	if got := w.sent.of(sandbox.NotifyRunFailed); len(got) != 1 {
		t.Fatalf("watchdog stop = %+v", got)
	}
}

func TestEarlyFailureIsRefundedAndLowQuotaWarns(t *testing.T) {
	w := newWorld(t)
	w.term()
	f := w.fleet(1)
	job := sandbox.RunInput{Kind: sandbox.RunJob, Title: "Обучение", Command: "python train.py", TemplateID: f.template.ID, ClassID: f.t4.ID, MaxHours: 14}
	broken, err := w.quotas.StartRun(w.ctx, w.aigerim, job)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	w.sync(f.node.Token)
	w.sync(f.node.Token, sandbox.RunReport{RunID: broken.ID, Phase: sandbox.PhaseRunning, Usage: busy("python")})
	w.exec(`UPDATE sandbox.runs SET started_at = now() - interval '3 minutes'`)
	w.sync(f.node.Token, sandbox.RunReport{RunID: broken.ID, Phase: sandbox.PhaseExited, ExitCode: 1, Log: "Traceback\nModuleNotFoundError: torchvision"})
	got, err := w.quotas.GetRun(w.ctx, w.aigerim, broken.ID)
	if err != nil || got.Status != sandbox.RunFailed || !got.Refunded || *got.Charged != 0 || !strings.Contains(got.Log, "torchvision") {
		t.Fatalf("refund = %+v, %v", got, err)
	}
	if runs, _ := w.quotas.Runs(w.ctx, w.aigerim); runs[0].Log != "" {
		t.Fatalf("list carries log")
	}
	if q := w.quota(w.aigerim); q.Hours != 30 {
		t.Fatalf("quota after refund = %+v", q)
	}

	long, err := w.quotas.StartRun(w.ctx, w.aigerim, job)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	w.sync(f.node.Token)
	w.sync(f.node.Token, sandbox.RunReport{RunID: long.ID, Phase: sandbox.PhaseRunning, Usage: &sandbox.RunUsage{CPU: 90, GPU: ptr(40.0)}})
	w.sync(f.node.Token, sandbox.RunReport{RunID: long.ID, Phase: sandbox.PhaseRunning, Usage: &sandbox.RunUsage{CPU: 90, GPU: ptr(20.0)}})
	w.exec(`UPDATE sandbox.runs SET started_at = now() - interval '13 hours 36 minutes' WHERE id = $1`, long.ID)
	w.sync(f.node.Token, sandbox.RunReport{RunID: long.ID, Phase: sandbox.PhaseExited})
	got, _ = w.quotas.GetRun(w.ctx, w.aigerim, long.ID)
	if got.GPULoad == nil || *got.GPULoad != 30 || *got.Charged != 27.2 {
		t.Fatalf("long = %+v", got)
	}
	low := w.sent.of(sandbox.NotifyQuotaLow)
	if len(low) != 1 || low[0].Params["hours"] != "2.8" || low[0].Params["total"] != "30" {
		t.Fatalf("low = %+v", low)
	}
}

func ptr(v float64) *float64 {
	return &v
}
