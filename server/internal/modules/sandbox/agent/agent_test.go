package agent_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/sandbox"
	"github.com/babakapaa/lodestar/server/internal/modules/sandbox/agent"
)

type fakeContainer struct {
	id     string
	name   string
	state  string
	exit   int
	labels map[string]string
	spec   map[string]any
	port   string
}

type fakeDocker struct {
	mu         sync.Mutex
	containers map[string]*fakeContainer
	pulled     []string
	next       int
}

func (d *fakeDocker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	path := strings.TrimPrefix(r.URL.Path, "/v1.43")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	switch {
	case r.Method == http.MethodGet && path == "/containers/json":
		var out []map[string]any
		for _, c := range d.containers {
			out = append(out, map[string]any{"Id": c.id, "State": c.state, "Labels": c.labels})
		}
		_ = json.NewEncoder(w).Encode(out)
	case r.Method == http.MethodGet && len(parts) == 3 && parts[2] == "json":
		c := d.containers[parts[1]]
		ports := map[string]any{}
		if c.port != "" {
			for key := range c.spec["ExposedPorts"].(map[string]any) {
				ports[key] = []map[string]string{{"HostIp": "0.0.0.0", "HostPort": c.port}}
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"State": map[string]any{"ExitCode": c.exit}, "NetworkSettings": map[string]any{"Ports": ports}})
	case r.Method == http.MethodPost && path == "/images/create":
		d.pulled = append(d.pulled, r.URL.Query().Get("fromImage"))
		_, _ = w.Write([]byte(`{"status":"done"}`))
	case r.Method == http.MethodPost && path == "/containers/create":
		var spec map[string]any
		_ = json.NewDecoder(r.Body).Decode(&spec)
		d.next++
		id := fmt.Sprintf("c%d", d.next)
		labels := map[string]string{}
		for k, v := range spec["Labels"].(map[string]any) {
			labels[k] = v.(string)
		}
		d.containers[id] = &fakeContainer{id: id, name: r.URL.Query().Get("name"), state: "created", labels: labels, spec: spec}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{"Id": id})
	case r.Method == http.MethodPost && len(parts) == 3 && parts[2] == "start":
		c := d.containers[parts[1]]
		c.state = "running"
		if _, ok := c.spec["ExposedPorts"]; ok {
			c.port = fmt.Sprintf("4915%d", d.next)
		}
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodPost && len(parts) == 3 && parts[2] == "stop":
		c := d.containers[parts[1]]
		c.state, c.exit = "exited", 137
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodDelete && len(parts) == 2:
		delete(d.containers, parts[1])
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "unexpected "+r.Method+" "+path, http.StatusNotFound)
	}
}

func (d *fakeDocker) byRun(runID string) *fakeContainer {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range d.containers {
		if c.labels["lodestar.run"] == runID {
			return c
		}
	}
	return nil
}

type fakeLodestar struct {
	mu      sync.Mutex
	plan    []sandbox.RunSpec
	reports []sandbox.AgentReport
}

func (f *fakeLodestar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if r.URL.Path != "/v1/open/sandbox/agent/sync" || r.Header.Get("Authorization") != "Bearer lsn_t_secret" {
		http.Error(w, "no", http.StatusUnauthorized)
		return
	}
	var report sandbox.AgentReport
	_ = json.NewDecoder(r.Body).Decode(&report)
	f.reports = append(f.reports, report)
	_ = json.NewEncoder(w).Encode(sandbox.AgentPlan{NodeID: "n1", Runs: f.plan})
}

func (f *fakeLodestar) set(plan ...sandbox.RunSpec) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.plan = plan
}

func (f *fakeLodestar) last() sandbox.AgentReport {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.reports[len(f.reports)-1]
}

type rig struct {
	t      *testing.T
	docker *fakeDocker
	server *fakeLodestar
	agent  agent.Agent
	data   string
}

func newRig(t *testing.T, gpus int) rig {
	t.Helper()
	docker := &fakeDocker{containers: map[string]*fakeContainer{}}
	dockerServer := httptest.NewServer(docker)
	t.Cleanup(dockerServer.Close)
	lodestar := &fakeLodestar{}
	lodestarServer := httptest.NewServer(lodestar)
	t.Cleanup(lodestarServer.Close)
	client, err := agent.NewDocker(dockerServer.URL)
	if err != nil {
		t.Fatalf("docker: %v", err)
	}
	data := t.TempDir()
	cfg := agent.Config{Server: lodestarServer.URL, Token: "lsn_t_secret", GPUs: gpus, Host: "gpu-01.kaznu.kz", DataDir: data}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("config: %v", err)
	}
	return rig{t: t, docker: docker, server: lodestar, agent: agent.New(cfg, client), data: data}
}

func (r rig) tick() {
	r.t.Helper()
	if err := r.agent.Tick(context.Background()); err != nil {
		r.t.Fatalf("tick: %v", err)
	}
}

func notebook(id string) sandbox.RunSpec {
	return sandbox.RunSpec{RunID: id, Kind: sandbox.RunJupyter, Title: "Ноутбук", Image: "pytorch/pytorch:2.5.1", MaxHours: 4, Workspace: "user-u1"}
}

func TestAgentStartsReportsAndStopsSession(t *testing.T) {
	r := newRig(t, 1)
	r.server.set(notebook("r1"))
	r.tick()
	c := r.docker.byRun("r1")
	if c == nil || c.state != "running" || c.name != "lodestar-r1" || len(r.docker.pulled) != 1 {
		t.Fatalf("container = %+v, pulled %v", c, r.docker.pulled)
	}
	host := c.spec["HostConfig"].(map[string]any)
	device := host["DeviceRequests"].([]any)[0].(map[string]any)
	if ids := device["DeviceIDs"].([]any); ids[0] != "0" || host["Binds"].([]any)[0] != r.data+"/user-u1:/workspace" {
		t.Fatalf("host config = %+v", host)
	}
	if cmd := fmt.Sprint(c.spec["Cmd"]); !strings.Contains(cmd, "jupyter lab") || !strings.Contains(cmd, "--ServerApp.token="+c.labels["lodestar.token"]) {
		t.Fatalf("cmd = %s", cmd)
	}

	r.tick()
	report := r.server.last()
	want := "http://gpu-01.kaznu.kz:" + c.port + "/lab?token=" + c.labels["lodestar.token"]
	if len(report.Runs) != 1 || report.Runs[0].Phase != sandbox.PhaseRunning || report.Runs[0].URL != want {
		t.Fatalf("report = %+v, want url %s", report, want)
	}

	stop := notebook("r1")
	stop.Stop = true
	r.server.set(stop)
	r.tick()
	if c := r.docker.byRun("r1"); c.state != "exited" {
		t.Fatalf("not stopped = %+v", c)
	}
	r.server.set()
	r.tick()
	if report := r.server.last(); report.Runs[0].Phase != sandbox.PhaseExited || report.Runs[0].ExitCode != 137 {
		t.Fatalf("exit report = %+v", report)
	}
	if c := r.docker.byRun("r1"); c != nil {
		t.Fatalf("container left = %+v", c)
	}
}

func TestAgentGivesEachRunItsOwnGPU(t *testing.T) {
	r := newRig(t, 2)
	job := sandbox.RunSpec{RunID: "j1", Kind: sandbox.RunJob, Title: "Обучение", Image: "img", Command: "python train.py", MaxHours: 2, Workspace: "team-t1"}
	other := job
	other.RunID = "j2"
	third := job
	third.RunID = "j3"
	r.server.set(job, other, third)
	r.tick()
	first, second := r.docker.byRun("j1"), r.docker.byRun("j2")
	if first == nil || second == nil || r.docker.byRun("j3") != nil {
		t.Fatalf("containers = %+v %+v", first, second)
	}
	if first.labels["lodestar.gpu"] == second.labels["lodestar.gpu"] {
		t.Fatalf("same gpu = %s", first.labels["lodestar.gpu"])
	}
	if cmd := fmt.Sprint(first.spec["Cmd"]); cmd != "[sh -c python train.py]" {
		t.Fatalf("cmd = %s", cmd)
	}
	if _, ok := first.spec["ExposedPorts"]; ok {
		t.Fatalf("job exposes port")
	}
	r.server.set(other, third)
	r.tick()
	if r.docker.byRun("j1") != nil || r.docker.byRun("j3") == nil {
		t.Fatalf("freed gpu not reused")
	}
}

func TestAgentRefusesStrangeWorkspace(t *testing.T) {
	r := newRig(t, 1)
	bad := notebook("r1")
	bad.Workspace = "../etc"
	r.server.set(bad)
	r.tick()
	if r.docker.byRun("r1") != nil {
		t.Fatalf("launched with bad workspace")
	}
	if err := (agent.Config{Server: "x", Token: "y", GPUs: 1, DataDir: "/d"}).Validate(); err == nil {
		t.Fatalf("config without host accepted")
	}
}
