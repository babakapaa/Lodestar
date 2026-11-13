package navigator

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type ModelReadiness struct {
	settings   module.Settings
	configured bool
}

func (r ModelReadiness) Problems(ctx context.Context, _ pgx.Tx, tenantID string) ([]module.Problem, error) {
	if r.configured {
		return nil, nil
	}
	if on, err := r.settings.ModuleEnabled(ctx, tenantID, Name); err != nil || !on {
		return nil, err
	}
	return []module.Problem{{Section: "navigator", Code: "model_missing", Count: 1}}, nil
}
