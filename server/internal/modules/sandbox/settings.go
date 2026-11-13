package sandbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const BaseHours = "sandbox.base_hours"

func (Module) Settings() []schema.Definition {
	least, most := 0, 10000
	return []schema.Definition{
		{Key: BaseHours, Type: schema.Integer, Min: &least, Max: &most, Default: 20},
	}
}

func (s Service) baseHours(ctx context.Context, tx pgx.Tx) (float64, error) {
	raw, err := s.settings.Value(ctx, tx, BaseHours)
	if err != nil {
		return 0, fmt.Errorf("read %s: %w", BaseHours, err)
	}
	var hours int
	if err := json.Unmarshal(raw, &hours); err != nil {
		return 0, fmt.Errorf("read %s: %w", BaseHours, err)
	}
	return float64(hours), nil
}
