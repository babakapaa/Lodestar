package sandbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	BaseHours    = "sandbox.base_hours"
	DeadlineDays = "sandbox.deadline_days"
)

func (Module) Settings() []schema.Definition {
	least, most, days := 0, 10000, 30
	return []schema.Definition{
		{Key: BaseHours, Type: schema.Integer, Min: &least, Max: &most, Default: 20},
		{Key: DeadlineDays, Type: schema.Integer, Min: &least, Max: &days, Default: 3},
	}
}

func (s Service) baseHours(ctx context.Context, tx pgx.Tx) (float64, error) {
	h, err := s.integer(ctx, tx, BaseHours)
	return float64(h), err
}

func (s Service) integer(ctx context.Context, tx pgx.Tx, key string) (int, error) {
	raw, err := s.settings.Value(ctx, tx, key)
	if err != nil {
		return 0, fmt.Errorf("read %s: %w", key, err)
	}
	var v int
	if err := json.Unmarshal(raw, &v); err != nil {
		return 0, fmt.Errorf("read %s: %w", key, err)
	}
	return v, nil
}
