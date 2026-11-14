package setup_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/setup"
)

type fakeProbe struct {
	name     string
	required bool
	err      error
}

func (p *fakeProbe) Name() string                { return p.name }
func (p *fakeProbe) Required() bool              { return p.required }
func (p *fakeProbe) Check(context.Context) error { return p.err }

type gpuNodes struct{}

func (gpuNodes) Hardware(context.Context, pgx.Tx, string) ([]module.HardwareCheck, error) {
	return []module.HardwareCheck{{Name: "gpu:a100-1", Status: module.CheckFailed}}, nil
}

func problemCounts(problems []module.Problem) map[string]int {
	out := map[string]int{}
	for _, p := range problems {
		out[p.Code] = p.Count
	}
	return out
}

func TestHardwareChecksGateLaunch(t *testing.T) {
	w := newWorld(t)
	if err := w.setup.Begin(w.ctx, w.tenant); err != nil {
		t.Fatalf("begin: %v", err)
	}
	storage := &fakeProbe{name: "storage", required: true, err: errors.New("bucket is gone")}
	w.setup = w.setup.WithHardware([]module.HardwareSource{gpuNodes{}},
		&fakeProbe{name: "database", required: true}, storage,
		&fakeProbe{name: "speech", err: setup.ErrNotConfigured})

	o, err := w.setup.Overview(w.ctx, w.admin)
	if err != nil || problemCounts(o.Problems)["hardware_unchecked"] != 1 {
		t.Fatalf("before check = %+v, %v", o.Problems, err)
	}
	report, err := w.setup.CheckHardware(w.ctx, w.admin)
	if err != nil || report.CheckedAt == nil || len(report.Checks) != 4 {
		t.Fatalf("check = %+v, %v", report, err)
	}
	byName := map[string]module.HardwareCheck{}
	for _, c := range report.Checks {
		byName[c.Name] = c
	}
	if byName["database"].Status != module.CheckOK || byName["storage"].Detail != "bucket is gone" || byName["speech"].Status != module.CheckOff || byName["gpu:a100-1"].Status != module.CheckFailed {
		t.Fatalf("checks = %+v", report.Checks)
	}
	o, _ = w.setup.Overview(w.ctx, w.admin)
	if got := problemCounts(o.Problems); got["hardware_failed"] != 1 || got["hardware_unchecked"] != 0 {
		t.Fatalf("after failure = %+v", o.Problems)
	}

	storage.err = nil
	if _, err := w.setup.CheckHardware(w.ctx, w.admin); err != nil {
		t.Fatalf("recheck: %v", err)
	}
	o, _ = w.setup.Overview(w.ctx, w.admin)
	if got := problemCounts(o.Problems); got["hardware_failed"] != 0 {
		t.Fatalf("after fix = %+v", o.Problems)
	}
	saved, err := w.setup.Hardware(w.ctx, w.admin)
	if err != nil || len(saved.Checks) != 4 {
		t.Fatalf("saved = %+v, %v", saved, err)
	}

	stranger := w.user()
	if _, err := w.setup.CheckHardware(w.ctx, stranger); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("stranger: %v", err)
	}
	if _, err := w.setup.Assign(w.ctx, w.admin, setup.SectionDeploy, setup.AssignInput{UserID: &stranger.UserID}); err != nil {
		t.Fatalf("assign: %v", err)
	}
	if _, err := w.setup.CheckHardware(w.ctx, stranger); err != nil {
		t.Fatalf("assignee checks: %v", err)
	}
}
