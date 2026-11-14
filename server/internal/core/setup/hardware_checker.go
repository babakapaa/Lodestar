package setup

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const probeTimeout = 10 * time.Second

type hardwareChecker struct {
	probes []Probe
}

func (c hardwareChecker) run(ctx context.Context) []module.HardwareCheck {
	out := make([]module.HardwareCheck, len(c.probes))
	var wg sync.WaitGroup
	for i, p := range c.probes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out[i] = c.one(ctx, p)
		}()
	}
	wg.Wait()
	return out
}

func (hardwareChecker) one(ctx context.Context, p Probe) module.HardwareCheck {
	ctx, cancel := context.WithTimeout(ctx, probeTimeout)
	defer cancel()
	started := time.Now()
	err := p.Check(ctx)
	check := module.HardwareCheck{Name: p.Name(), Required: p.Required(), Status: module.CheckOK, Millis: time.Since(started).Milliseconds()}
	switch {
	case errors.Is(err, ErrNotConfigured):
		check.Status = module.CheckOff
	case err != nil:
		check.Status, check.Detail = module.CheckFailed, err.Error()
	}
	return check
}
