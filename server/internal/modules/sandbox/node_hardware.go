package sandbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type NodeHardware struct {
	settings module.Settings
	nodes    NodeStore
}

func (h NodeHardware) Hardware(ctx context.Context, tx pgx.Tx, tenantID string) ([]module.HardwareCheck, error) {
	if on, err := h.settings.ModuleEnabled(ctx, tenantID, Name); err != nil || !on {
		return nil, err
	}
	nodes, err := h.nodes.All(ctx, tx)
	if err != nil {
		return nil, err
	}
	out := make([]module.HardwareCheck, 0, len(nodes))
	for _, n := range nodes {
		check := module.HardwareCheck{Name: "gpu:" + n.Name, Status: module.CheckOK, Detail: n.ClassName}
		if !n.Online {
			check.Status = module.CheckFailed
		}
		out = append(out, check)
	}
	return out, nil
}
