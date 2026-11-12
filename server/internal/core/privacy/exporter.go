package privacy

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Exporter struct {
	tenants  database.TenantTx
	sections Sections
}

func NewExporter(tenants database.TenantTx, sections Sections) Exporter {
	return Exporter{tenants: tenants, sections: sections}
}

func (e Exporter) Collect(ctx context.Context, p auth.Principal) (Archive, error) {
	data := map[string]any{}
	var blobs []Blob
	err := e.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		for _, s := range e.sections {
			v, err := s.Export(ctx, tx, p.UserID)
			if err != nil {
				return fmt.Errorf("export %s: %w", s.Name(), err)
			}
			data[s.Name()] = v
			if source, ok := s.(BlobSource); ok {
				found, err := source.Blobs(ctx, tx, p.TenantID, p.UserID)
				if err != nil {
					return fmt.Errorf("export %s files: %w", s.Name(), err)
				}
				blobs = append(blobs, found...)
			}
		}
		return nil
	})
	if err != nil {
		return Archive{}, err
	}
	return Archive{userID: p.UserID, sections: e.sections, data: data, blobs: blobs}, nil
}
