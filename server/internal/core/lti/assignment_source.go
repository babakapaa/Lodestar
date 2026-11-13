package lti

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/schema"
)

const (
	SourceSetting  = "core.lti_assignment_source"
	SourceLodestar = "lodestar"
	SourceLMS      = "lms"
	SourceUnits    = "units"
)

type Values interface {
	Value(ctx context.Context, tx pgx.Tx, key string) (json.RawMessage, error)
}

func CoreSettings() []schema.Definition {
	return []schema.Definition{
		{Key: SourceSetting, Type: schema.Choice, Options: []string{SourceLodestar, SourceLMS, SourceUnits}, Default: SourceLodestar},
	}
}

func tenantSource(ctx context.Context, tx pgx.Tx, values Values) (string, error) {
	raw, err := values.Value(ctx, tx, SourceSetting)
	if err != nil {
		return "", err
	}
	var v string
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", fmt.Errorf("read %s: %w", SourceSetting, err)
	}
	return v, nil
}

func effectiveSource(tenant string, c Context) string {
	switch {
	case tenant != SourceUnits:
		return tenant
	case c.Source != nil:
		return *c.Source
	default:
		return SourceLodestar
	}
}
