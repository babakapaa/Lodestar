package sandbox_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/sandbox"
)

func TestEstimateSuggestsClassAndShowsNight(t *testing.T) {
	w := newWorld(t)
	w.term()
	f := w.fleet(1)
	a100, err := w.quotas.CreateClass(w.ctx, w.admin, sandbox.GPUClassInput{Name: "A100", MemoryGB: 80, Rate: 8})
	if err != nil {
		t.Fatalf("a100: %v", err)
	}
	if _, err := w.quotas.Estimate(w.ctx, w.aigerim, sandbox.EstimateInput{Kind: "train", TemplateID: f.template.ID}); !errors.Is(err, sandbox.ErrInvalidRun) {
		t.Fatalf("bad kind: %v", err)
	}
	w.sync(f.node.Token)
	e, err := w.quotas.Estimate(w.ctx, w.aigerim, sandbox.EstimateInput{Kind: sandbox.RunJob, TemplateID: f.template.ID, MaxHours: 4})
	if err != nil || e.Available != 30 || e.Night != nil || len(e.Options) != 2 || e.Suggested == nil || *e.Suggested != f.t4.ID {
		t.Fatalf("estimate = %+v, %v", e, err)
	}
	for _, o := range e.Options {
		switch o.ClassID {
		case f.t4.ID:
			if o.Cost != 8 || o.FreeGPUs != 1 || !o.Enough || o.NightCost != 8 {
				t.Fatalf("t4 = %+v", o)
			}
		case a100.ID:
			if o.Cost != 32 || o.FreeGPUs != 0 || o.Enough {
				t.Fatalf("a100 = %+v", o)
			}
		}
	}
	w.settings.set(sandbox.NightPercent, 50)
	w.settings.set(sandbox.NightFrom, 23)
	w.settings.set(sandbox.NightTo, 7)
	w.settings.set(sandbox.TimeZone, "Asia/Almaty")
	if _, err := w.quotas.StartRun(w.ctx, w.erlan, sandbox.RunInput{Kind: sandbox.RunJob, Title: "x", Command: "x", TemplateID: f.template.ID, ClassID: f.t4.ID, MaxHours: 1}); err != nil {
		t.Fatalf("erlan: %v", err)
	}
	if _, err := w.quotas.StartRun(w.ctx, w.erlan, sandbox.RunInput{Kind: sandbox.RunJob, Title: "y", Command: "y", TemplateID: f.template.ID, ClassID: f.t4.ID, MaxHours: 1}); err != nil {
		t.Fatalf("erlan: %v", err)
	}
	rec := w.call(w.aigerim, http.MethodPost, "/v1/sandbox/estimate", `{"kind":"jupyter","templateId":"`+f.template.ID+`"}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"timeZone":"Asia/Almaty"`) || !strings.Contains(rec.Body.String(), `"nightCost":4`) {
		t.Fatalf("night estimate = %d %s", rec.Code, rec.Body)
	}
}

func TestStarterDemoRerunExportAndPace(t *testing.T) {
	w := newWorld(t)
	w.term()
	f := w.fleet(2)
	starter, err := w.quotas.CreateTemplate(w.ctx, w.admin, sandbox.TemplateInput{Name: "Первая нейросеть", Image: "lodestar/starter:1", Starter: true})
	if err != nil || !starter.Starter {
		t.Fatalf("starter template = %+v, %v", starter, err)
	}
	free := sandbox.RunInput{Kind: sandbox.RunJupyter, Title: "Первый шаг", TemplateID: starter.ID, ClassID: f.t4.ID}
	run, err := w.quotas.StartRun(w.ctx, w.aigerim, free)
	if err != nil || run.Rate != 0 || run.MaxHours != 0.25 {
		t.Fatalf("starter run = %+v, %v", run, err)
	}
	if _, err := w.quotas.StartRun(w.ctx, w.aigerim, free); !errors.Is(err, sandbox.ErrStarterBusy) {
		t.Fatalf("second starter: %v", err)
	}
	team := free
	team.TeamID = "team-1"
	if _, err := w.quotas.StartRun(w.ctx, w.erlan, team); !errors.Is(err, sandbox.ErrInvalidRun) {
		t.Fatalf("team starter: %v", err)
	}
	if q := w.quota(w.aigerim); q.Reserved != 0 {
		t.Fatalf("starter reserves = %+v", q)
	}

	if _, err := w.quotas.StartRun(w.ctx, w.aigerim, sandbox.RunInput{Kind: sandbox.RunDemo, Title: "Демо", TemplateID: f.template.ID, ClassID: f.t4.ID}); !errors.Is(err, sandbox.ErrInvalidRun) {
		t.Fatalf("demo without command: %v", err)
	}
	demo, err := w.quotas.StartRun(w.ctx, w.aigerim, sandbox.RunInput{Kind: sandbox.RunDemo, Title: "Демо", Command: "python app.py", TemplateID: f.template.ID, ClassID: f.t4.ID})
	if err != nil || demo.MaxHours != 2 {
		t.Fatalf("demo = %+v, %v", demo, err)
	}
	if _, err := w.quotas.RerunRun(w.ctx, w.aigerim, demo.ID); !errors.Is(err, sandbox.ErrRunActive) {
		t.Fatalf("rerun active: %v", err)
	}
	w.sync(f.node.Token)
	w.sync(f.node.Token,
		sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseRunning},
		sandbox.RunReport{RunID: demo.ID, Phase: sandbox.PhaseRunning, URL: "http://gpu:49153/"})
	w.exec(`UPDATE sandbox.runs SET started_at = now() - interval '90 minutes'`)
	w.sync(f.node.Token,
		sandbox.RunReport{RunID: run.ID, Phase: sandbox.PhaseExited},
		sandbox.RunReport{RunID: demo.ID, Phase: sandbox.PhaseExited})
	ended, _ := w.quotas.GetRun(w.ctx, w.aigerim, run.ID)
	if ended.Status != sandbox.RunFinished || *ended.Charged != 0 {
		t.Fatalf("starter charged = %+v", ended)
	}
	q := w.quota(w.aigerim)
	if q.Spent != 3 || q.Hours != 27 || q.Forecast >= q.Hours {
		t.Fatalf("pace = %+v", q)
	}

	if _, err := w.quotas.RerunRun(w.ctx, w.erlan, demo.ID); !errors.Is(err, sandbox.ErrRunNotFound) {
		t.Fatalf("stranger reruns: %v", err)
	}
	rec := w.call(w.aigerim, http.MethodPost, "/v1/sandbox/runs/"+demo.ID+"/rerun", "")
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"command":"python app.py"`) || !strings.Contains(rec.Body.String(), `"status":"queued"`) {
		t.Fatalf("rerun = %d %s", rec.Code, rec.Body)
	}

	rec = w.call(w.aigerim, http.MethodGet, "/v1/sandbox/runs/export?format=csv", "")
	body := rec.Body.String()
	if rec.Code != http.StatusOK || !strings.Contains(body, "queued_at,started_at") || strings.Count(body, "Демо") != 2 || !strings.Contains(body, "Первый шаг") {
		t.Fatalf("export = %d %s", rec.Code, body)
	}
	if rec := w.call(w.aigerim, http.MethodGet, "/v1/sandbox/runs/export?format=pdf", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("pdf export = %d", rec.Code)
	}
}
