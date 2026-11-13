package sandbox

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	BaseHours    = "sandbox.base_hours"
	DeadlineDays = "sandbox.deadline_days"
	Watch        = "sandbox.watch_processes"
)

const defaultWatch = "xmrig, xmr-stak, ethminer, t-rex, nbminer, lolminer, phoenixminer, gminer, teamredminer, nanominer, cpuminer, minerd, bzminer, rigel, srbminer"

func (Module) Settings() []schema.Definition {
	least, most, days := 0, 10000, 30
	return []schema.Definition{
		{Key: BaseHours, Type: schema.Integer, Min: &least, Max: &most, Default: 20},
		{Key: DeadlineDays, Type: schema.Integer, Min: &least, Max: &days, Default: 3},
		{Key: Watch, Type: schema.Text, MaxLength: 2000, Default: defaultWatch},
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

func (s Service) watched(ctx context.Context, tx pgx.Tx) ([]string, error) {
	raw, err := s.settings.Value(ctx, tx, Watch)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", Watch, err)
	}
	var v *string
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("read %s: %w", Watch, err)
	}
	list := defaultWatch
	if v != nil {
		list = *v
	}
	var out []string
	for _, name := range strings.Split(list, ",") {
		if name = strings.ToLower(strings.TrimSpace(name)); name != "" {
			out = append(out, name)
		}
	}
	return out, nil
}
