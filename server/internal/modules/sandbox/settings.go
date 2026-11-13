package sandbox

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	BaseHours    = "sandbox.base_hours"
	DeadlineDays = "sandbox.deadline_days"
	Watch        = "sandbox.watch_processes"
	NightPercent = "sandbox.night_discount"
	NightFrom    = "sandbox.night_from"
	NightTo      = "sandbox.night_to"
	TimeZone     = "core.time_zone"
)

const defaultWatch = "xmrig, xmr-stak, ethminer, t-rex, nbminer, lolminer, phoenixminer, gminer, teamredminer, nanominer, cpuminer, minerd, bzminer, rigel, srbminer"

func (Module) Settings() []schema.Definition {
	least, most, days, percent, hour := 0, 10000, 30, 90, 23
	return []schema.Definition{
		{Key: BaseHours, Type: schema.Integer, Min: &least, Max: &most, Default: 20},
		{Key: DeadlineDays, Type: schema.Integer, Min: &least, Max: &days, Default: 3},
		{Key: Watch, Type: schema.Text, MaxLength: 2000, Default: defaultWatch},
		{Key: NightPercent, Type: schema.Integer, Min: &least, Max: &percent, Default: 50},
		{Key: NightFrom, Type: schema.Integer, Min: &least, Max: &hour, Default: 23},
		{Key: NightTo, Type: schema.Integer, Min: &least, Max: &hour, Default: 7},
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

func (s Service) tariff(ctx context.Context, tx pgx.Tx) (tariff, error) {
	var t tariff
	var err error
	if t.discount, err = s.integer(ctx, tx, NightPercent); err != nil {
		return t, err
	}
	if t.from, err = s.integer(ctx, tx, NightFrom); err != nil {
		return t, err
	}
	if t.to, err = s.integer(ctx, tx, NightTo); err != nil {
		return t, err
	}
	t.loc, _ = s.location(ctx, tx)
	return t, nil
}

func (s Service) location(ctx context.Context, tx pgx.Tx) (*time.Location, string) {
	raw, err := s.settings.Value(ctx, tx, TimeZone)
	if err != nil {
		return time.UTC, "UTC"
	}
	var name string
	if json.Unmarshal(raw, &name) != nil || name == "" {
		return time.UTC, "UTC"
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC, "UTC"
	}
	return loc, name
}
