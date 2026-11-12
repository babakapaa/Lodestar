package notes

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type SourcesPrivacySection struct {
	service Service
}

func (SourcesPrivacySection) Name() string {
	return "notes.sources"
}

func (s SourcesPrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	rows, err := s.service.sources.Owned(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]Source, 0, len(rows))
	for _, row := range rows {
		src, err := s.service.source(row)
		if err != nil {
			return nil, fmt.Errorf("export source: %w", err)
		}
		out = append(out, src)
	}
	return out, nil
}
