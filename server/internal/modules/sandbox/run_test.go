package sandbox_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/sandbox"
)

func TestRunGoesThroughQueueNodeAndQuota(t *testing.T) {
	w := newWorld(t)
	w.term()
	f := w.fleet(1)
	if _, err := w.quotas.CreateNode(w.ctx, w.teacher, sandbox.NodeInput{Name: "x", ClassID: f.t4.ID, GPUs: 1}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher adds node: %v", err)
	}
	if !strings.HasPrefix(f.node.Token, "lsn_"+w.tenant+"_") || f.node.Online {
		t.Fatalf("node = %+v", f.node)
	}
	in := sandbox.RunInput{Kind: sandbox.RunJupyter, Title: "Ноутбук", TemplateID: f.template.ID, ClassID: f.t4.ID}
	if _, err := w.quotas.StartRun(w.ctx, w.aigerim, sandbox.RunInput{Kind: sandbox.RunJob, Title: "x", TemplateID: f.template.ID, ClassID: f.t4.ID}); !errors.Is(err, sandbox.ErrInvalidRun) {
		t.Fatalf("job without command: %v", err)
	}
	big := in
	big.MaxHours = 100
	if _, err := w.quotas.StartRun(w.ctx, w.aigerim, big); !errors.Is(err, sandbox.ErrNotEnoughHours) {
		t.Fatalf("too long: %v", err)
	}
	run, err := w.quotas.StartRun(w.ctx, w.aigerim, in)
	if err != nil || run.Status != sandbox.RunQueued || run.MaxHours != 4 || run.Position == nil || *run.Position != 1 {
		t.Fatalf("start = %+v, %v", run, err)
	}
	if q := w.quota(w.aigerim); q.Hours != 30 || q.Reserved != 8 {
		t.Fatalf("quota while queued = %+v", q)
	}

	plan := w.sync(f.node.Token)
	if len(plan.Runs) != 1 || plan.Runs[0].RunID != run.ID || plan.Runs[0].Image != f.template.Image || plan.Runs[0].Stop ||
		plan.Runs[0].Workspace != "user-"+w.aigerim.UserID {
		t.Fatalf("plan = %+v", plan)
	}
	w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseRunning, URL: "https://gpu-01.kaznu.kz/s/abc"})
	mine, err := w.quotas.GetRun(w.ctx, w.aigerim, run.ID)
	if err != nil || mine.Status != sandbox.RunRunning || mine.URL == "" || mine.StartedAt == nil {
		t.Fatalf("running = %+v, %v", mine, err)
	}
	if _, err := w.quotas.GetRun(w.ctx, w.erlan, run.ID); !errors.Is(err, sandbox.ErrRunNotFound) {
		t.Fatalf("stranger sees run: %v", err)
	}
	w.exec(`UPDATE sandbox.runs SET started_at = now() - interval '90 minutes' WHERE id = $1`, run.ID)

	if _, err := w.quotas.StopRun(w.ctx, w.erlan, run.ID); !errors.Is(err, sandbox.ErrRunNotFound) {
		t.Fatalf("stranger stops: %v", err)
	}
	stopping, err := w.quotas.StopRun(w.ctx, w.aigerim, run.ID)
	if err != nil || stopping.Status != sandbox.RunStopping {
		t.Fatalf("stop = %+v, %v", stopping, err)
	}
	if plan := w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseRunning}); !plan.Runs[0].Stop {
		t.Fatalf("plan keeps run = %+v", plan)
	}
	w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseExited, ExitCode: 137})
	ended, err := w.quotas.GetRun(w.ctx, w.aigerim, run.ID)
	if err != nil || ended.Status != sandbox.RunStopped || ended.Charged == nil || *ended.Charged != 3 || ended.URL != "" {
		t.Fatalf("ended = %+v, %v", ended, err)
	}
	q := w.quota(w.aigerim)
	if q.Hours != 27 || q.Reserved != 0 || q.Lines[len(q.Lines)-1].Source != sandbox.SourceRun {
		t.Fatalf("quota after run = %+v", q)
	}
	if _, err := w.quotas.StopRun(w.ctx, w.aigerim, run.ID); !errors.Is(err, sandbox.ErrRunEnded) {
		t.Fatalf("stop twice: %v", err)
	}
	if nodes, err := w.quotas.Nodes(w.ctx, w.admin); err != nil || !nodes[0].Online || nodes[0].Busy != 0 {
		t.Fatalf("nodes = %+v, %v", nodes, err)
	}
	var metrics []module.Metric
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		metrics, err = sandbox.SandboxStats{}.Stats(w.ctx, tx, []string{w.aigerim.UserID, w.erlan.UserID}, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
		return err
	})
	if err != nil || len(metrics) != 1 || metrics[0].Key != sandbox.MetricHours || metrics[0].Value != 3 {
		t.Fatalf("stats = %+v, %v", metrics, err)
	}
}

func TestCourseWorkNearDeadlineGoesFirst(t *testing.T) {
	w := newWorld(t)
	w.term()
	f := w.fleet(1)
	project := sandbox.RunInput{Kind: sandbox.RunJob, Title: "Свой проект", Command: "python train.py", TemplateID: f.template.ID, ClassID: f.t4.ID, MaxHours: 2}
	if _, err := w.quotas.StartRun(w.ctx, w.erlan, project); err != nil {
		t.Fatalf("erlan project: %v", err)
	}
	second, err := w.quotas.StartRun(w.ctx, w.aigerim, project)
	if err != nil {
		t.Fatalf("aigerim project: %v", err)
	}
	course := project
	course.Title = "Лабораторная"
	course.WorkKind, course.WorkID = "assignment", "a1"
	homework, err := w.quotas.StartRun(w.ctx, w.aigerim, course)
	if err != nil || homework.WorkTitle != "Нейросеть для МРТ" || *homework.Position != 1 {
		t.Fatalf("homework = %+v, %v", homework, err)
	}
	runs, err := w.quotas.Runs(w.ctx, w.erlan)
	if err != nil || len(runs) != 1 || *runs[0].Position != 3 {
		t.Fatalf("erlan waits = %+v, %v", runs, err)
	}
	if mine, _ := w.quotas.GetRun(w.ctx, w.aigerim, second.ID); *mine.Position != 2 {
		t.Fatalf("division goes before = %+v", mine)
	}
	plan := w.sync(f.node.Token)
	if len(plan.Runs) != 1 || plan.Runs[0].RunID != homework.ID {
		t.Fatalf("first on gpu = %+v", plan)
	}
}

func TestTeamPoolAndLostNode(t *testing.T) {
	w := newWorld(t)
	w.term()
	f := w.fleet(2)
	team := sandbox.RunInput{Kind: sandbox.RunVSCode, Title: "Код команды", TemplateID: f.template.ID, ClassID: f.t4.ID, MaxHours: 2, TeamID: "team-1"}
	if _, err := w.quotas.StartRun(w.ctx, w.erlan, team); !errors.Is(err, sandbox.ErrNotEnoughHours) {
		t.Fatalf("empty pool: %v", err)
	}
	if _, err := w.quotas.Transfer(w.ctx, w.aigerim, sandbox.TransferInput{TeamID: "team-1", Hours: 4}); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	run, err := w.quotas.StartRun(w.ctx, w.erlan, team)
	if err != nil || run.TeamName != "Нейроны" {
		t.Fatalf("team run = %+v, %v", run, err)
	}
	if _, err := w.quotas.Transfer(w.ctx, w.aigerim, sandbox.TransferInput{TeamID: "team-1", Hours: 999}); !errors.Is(err, sandbox.ErrNotEnoughHours) {
		t.Fatalf("over transfer: %v", err)
	}
	if runs, err := w.quotas.Runs(w.ctx, w.aigerim); err != nil || len(runs) != 1 || runs[0].URL != "" {
		t.Fatalf("teammate sees = %+v, %v", runs, err)
	}
	w.sync(f.node.Token)
	w.sync(f.node.Token, sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseRunning, URL: "https://gpu/s"})
	w.exec(`UPDATE sandbox.runs SET started_at = now() - interval '60 minutes' WHERE id = $1`, run.ID)
	w.exec(`UPDATE sandbox.nodes SET seen_at = now() - interval '30 minutes' WHERE id = $1`, f.node.ID)
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error { return w.quotas.Schedule(w.ctx, tx) })
	if err != nil {
		t.Fatalf("schedule: %v", err)
	}
	lost, err := w.quotas.GetRun(w.ctx, w.erlan, run.ID)
	if err != nil || lost.Status != sandbox.RunFailed || lost.Reason != sandbox.ReasonNodeLost || *lost.Charged != 1 {
		t.Fatalf("lost = %+v, %v", lost, err)
	}
	if q := w.quota(w.erlan); q.Teams[0].Hours != 3 || q.Teams[0].Reserved != 0 {
		t.Fatalf("pool after run = %+v", q.Teams)
	}
}

func TestAgentNeedsItsToken(t *testing.T) {
	w := newWorld(t)
	f := w.fleet(1)
	if _, err := w.quotas.Sync(w.ctx, "lsn_"+w.tenant+"_wrong", sandbox.AgentReport{}); !errors.Is(err, sandbox.ErrBadNodeToken) {
		t.Fatalf("wrong secret: %v", err)
	}
	if rec := w.open(http.MethodPost, "/v1/open/sandbox/agent/sync", "garbage", `{}`); rec.Code != http.StatusUnauthorized {
		t.Fatalf("garbage = %d %s", rec.Code, rec.Body)
	}
	if rec := w.open(http.MethodPost, "/v1/open/sandbox/agent/sync", f.node.Token, `{"runs":[]}`); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), f.node.ID) {
		t.Fatalf("sync = %d %s", rec.Code, rec.Body)
	}
	if err := w.quotas.DeleteNode(w.ctx, w.admin, f.node.ID); err != nil {
		t.Fatalf("delete node: %v", err)
	}
	if _, err := w.quotas.Sync(w.ctx, f.node.Token, sandbox.AgentReport{}); !errors.Is(err, sandbox.ErrBadNodeToken) {
		t.Fatalf("deleted node syncs: %v", err)
	}
}
