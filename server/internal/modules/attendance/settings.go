package attendance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/platform/clientip"
)

const (
	QR          = "attendance.qr"
	Code        = "attendance.code"
	Networks    = "attendance.networks"
	LateMinutes = "attendance.late_minutes"
	AbsentAlert = "attendance.absent_alert"
	TimeZone    = "core.time_zone"
)

func (Module) Settings() []schema.Definition {
	none, late, alert, most := 0, 10, 1, 60
	return []schema.Definition{
		{Key: QR, Type: schema.Bool, Default: true},
		{Key: Code, Type: schema.Bool, Default: true},
		{Key: Networks, Type: schema.Text, MaxLength: 2000, Default: ""},
		{Key: LateMinutes, Type: schema.Integer, Min: &none, Max: &most, Default: late},
		{Key: AbsentAlert, Type: schema.Integer, Min: &alert, Max: &most, Default: 3},
	}
}

func (s Service) flag(ctx context.Context, tx pgx.Tx, key string) (bool, error) {
	var v bool
	return v, s.read(ctx, tx, key, &v)
}

func (s Service) integer(ctx context.Context, tx pgx.Tx, key string) (int, error) {
	var v int
	return v, s.read(ctx, tx, key, &v)
}

func (s Service) methods(ctx context.Context, tx pgx.Tx) (Methods, error) {
	var m Methods
	var err error
	if m.QR, err = s.flag(ctx, tx, QR); err != nil {
		return m, err
	}
	m.Code, err = s.flag(ctx, tx, Code)
	return m, err
}

func (s Service) networks(ctx context.Context, tx pgx.Tx) ([]netip.Prefix, error) {
	var raw *string
	if err := s.read(ctx, tx, Networks, &raw); err != nil || raw == nil {
		return nil, err
	}
	var out []netip.Prefix
	for _, item := range strings.FieldsFunc(*raw, func(r rune) bool { return r == ',' || r == '\n' || r == ' ' }) {
		if prefix, err := clientip.ParsePrefix(item); err == nil {
			out = append(out, prefix)
		}
	}
	return out, nil
}

func (s Service) location(ctx context.Context, tx pgx.Tx) *time.Location {
	var name string
	if s.read(ctx, tx, TimeZone, &name) != nil || name == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return loc
}

func (s Service) read(ctx context.Context, tx pgx.Tx, key string, into any) error {
	raw, err := s.settings.Value(ctx, tx, key)
	if err != nil {
		return fmt.Errorf("read %s: %w", key, err)
	}
	if err := json.Unmarshal(raw, into); err != nil {
		return fmt.Errorf("read %s: %w", key, err)
	}
	return nil
}
