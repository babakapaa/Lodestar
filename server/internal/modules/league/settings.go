package league

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const DivisionGPUHours = "league.division_gpu_hours"

func (Module) Settings() []schema.Definition {
	least, most := 0, 10000
	return []schema.Definition{
		{Key: DivisionGPUHours, Type: schema.Integer, Min: &least, Max: &most, Default: 0},
	}
}

func (s Service) gpuCap(ctx context.Context, tx pgx.Tx) (int, error) {
	raw, err := s.settings.Value(ctx, tx, DivisionGPUHours)
	if err != nil {
		return 0, fmt.Errorf("read %s: %w", DivisionGPUHours, err)
	}
	var hours int
	if err := json.Unmarshal(raw, &hours); err != nil {
		return 0, fmt.Errorf("read %s: %w", DivisionGPUHours, err)
	}
	return hours, nil
}
