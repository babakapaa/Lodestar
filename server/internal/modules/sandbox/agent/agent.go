package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/babakapaa/lodestar/server/internal/modules/sandbox"
)

var workspacePattern = regexp.MustCompile(`^(user|team)-[A-Za-z0-9-]{1,100}$`)

const logLines = 100

type Agent struct {
	cfg    Config
	docker Docker
	meter  GPUMeter
	client *http.Client
}

func New(cfg Config, docker Docker, meter GPUMeter) Agent {
	if cfg.Scheme == "" {
		cfg.Scheme = "http"
	}
	if cfg.Every == 0 {
		cfg.Every = 10 * time.Second
	}
	return Agent{cfg: cfg, docker: docker, meter: meter, client: &http.Client{Timeout: 30 * time.Second}}
}

func (a Agent) Run(ctx context.Context) {
	ticker := time.NewTicker(a.cfg.Every)
	defer ticker.Stop()
	for {
		if err := a.Tick(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("sync with lodestar", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (a Agent) Tick(ctx context.Context) error {
	containers, err := a.docker.List(ctx)
	if err != nil {
		return err
	}
	plan, err := a.sync(ctx, a.report(ctx, containers))
	if err != nil {
		return err
	}
	return a.apply(ctx, containers, plan)
}

func (a Agent) report(ctx context.Context, containers []Container) sandbox.AgentReport {
	report := sandbox.AgentReport{Runs: []sandbox.RunReport{}}
	load := a.load(ctx)
	for _, c := range containers {
		switch {
		case c.running():
			report.Runs = append(report.Runs, sandbox.RunReport{
				RunID: c.RunID, Phase: sandbox.PhaseRunning, URL: link(a.cfg.Scheme, a.cfg.Host, c), Usage: a.usage(ctx, c, load),
			})
		case c.paused():
			report.Runs = append(report.Runs, sandbox.RunReport{RunID: c.RunID, Phase: sandbox.PhaseRunning, URL: link(a.cfg.Scheme, a.cfg.Host, c)})
		case c.exited():
			logs, err := a.docker.Logs(ctx, c.ID, logLines)
			if err != nil {
				slog.Warn("read run log", "run", c.RunID, "error", err)
			}
			report.Runs = append(report.Runs, sandbox.RunReport{RunID: c.RunID, Phase: sandbox.PhaseExited, ExitCode: c.ExitCode, Log: logs})
		}
	}
	return report
}

func (a Agent) load(ctx context.Context) map[int]float64 {
	if a.meter == nil {
		return nil
	}
	load, err := a.meter.Load(ctx)
	if err != nil {
		slog.Warn("read gpu load", "error", err)
		return nil
	}
	return load
}

func (a Agent) usage(ctx context.Context, c Container, load map[int]float64) *sandbox.RunUsage {
	s, err := a.docker.Stats(ctx, c.ID)
	if err != nil {
		slog.Warn("read run stats", "run", c.RunID, "error", err)
		return nil
	}
	processes, err := a.docker.Processes(ctx, c.ID)
	if err != nil {
		slog.Warn("read run processes", "run", c.RunID, "error", err)
		return nil
	}
	u := &sandbox.RunUsage{CPU: s.CPU, NetBytes: s.NetBytes, Processes: processes}
	if v, ok := load[c.GPU]; ok {
		u.GPU = &v
	}
	return u
}

func (a Agent) sync(ctx context.Context, report sandbox.AgentReport) (sandbox.AgentPlan, error) {
	raw, err := json.Marshal(report)
	if err != nil {
		return sandbox.AgentPlan{}, err
	}
	endpoint := strings.TrimSuffix(a.cfg.Server, "/") + "/v1/open/sandbox/agent/sync"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return sandbox.AgentPlan{}, err
	}
	req.Header.Set("Authorization", "Bearer "+a.cfg.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return sandbox.AgentPlan{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return sandbox.AgentPlan{}, fmt.Errorf("lodestar answered %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var plan sandbox.AgentPlan
	err = json.NewDecoder(resp.Body).Decode(&plan)
	return plan, err
}

func (a Agent) apply(ctx context.Context, containers []Container, plan sandbox.AgentPlan) error {
	wanted := make(map[string]sandbox.RunSpec, len(plan.Runs))
	for _, spec := range plan.Runs {
		wanted[spec.RunID] = spec
	}
	have := map[string]bool{}
	used := map[int]bool{}
	for _, c := range containers {
		spec, ok := wanted[c.RunID]
		switch {
		case !ok:
			if err := a.docker.Remove(ctx, c.ID); err != nil {
				return err
			}
			continue
		case spec.Stop && c.paused():
			if err := a.docker.Unpause(ctx, c.ID); err != nil {
				return err
			}
			if err := a.docker.Stop(ctx, c.ID); err != nil {
				return err
			}
		case spec.Stop && c.running():
			if err := a.docker.Stop(ctx, c.ID); err != nil {
				return err
			}
		case spec.Hold && c.running():
			if err := a.docker.Pause(ctx, c.ID); err != nil {
				return err
			}
		case !spec.Hold && c.paused():
			if err := a.docker.Unpause(ctx, c.ID); err != nil {
				return err
			}
		case !spec.Stop && c.State == "created":
			if err := a.docker.Start(ctx, c.ID); err != nil {
				slog.Warn("start run", "run", c.RunID, "error", err)
			}
		}
		have[c.RunID] = true
		if !c.exited() {
			used[c.GPU] = true
		}
	}
	for _, spec := range plan.Runs {
		if have[spec.RunID] || spec.Stop || spec.Hold {
			continue
		}
		gpu := a.freeGPU(used)
		if gpu < 0 {
			slog.Warn("no free gpu for run", "run", spec.RunID)
			continue
		}
		if err := a.launch(ctx, spec, gpu); err != nil {
			slog.Warn("launch run", "run", spec.RunID, "error", err)
			continue
		}
		used[gpu] = true
	}
	return nil
}

func (a Agent) freeGPU(used map[int]bool) int {
	for i := range a.cfg.GPUs {
		if !used[i] {
			return i
		}
	}
	return -1
}

func (a Agent) launch(ctx context.Context, spec sandbox.RunSpec, gpu int) error {
	if !workspacePattern.MatchString(spec.Workspace) {
		return fmt.Errorf("bad workspace %q", spec.Workspace)
	}
	dir := filepath.Join(a.cfg.DataDir, spec.Workspace)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	token, err := newToken()
	if err != nil {
		return err
	}
	if err := a.docker.Pull(ctx, spec.Image); err != nil {
		return err
	}
	id, err := a.docker.Create(ctx, "lodestar-"+spec.RunID, containerSpec(spec, gpu, token, dir))
	if err != nil {
		return err
	}
	return a.docker.Start(ctx, id)
}
