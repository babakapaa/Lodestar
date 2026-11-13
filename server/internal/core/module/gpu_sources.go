package module

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type GPUSources []GPUSource

func (s GPUSources) Grants(ctx context.Context, tx pgx.Tx, userID string, from, to time.Time) ([]GPUGrant, error) {
	var out []GPUGrant
	for _, source := range s {
		grants, err := source.GPUGrants(ctx, tx, userID, from, to)
		if err != nil {
			return nil, err
		}
		out = append(out, grants...)
	}
	return out, nil
}
