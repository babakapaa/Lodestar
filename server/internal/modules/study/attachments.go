package study

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/files"
)

type Attachments struct {
	files files.FileStore
}

func (a Attachments) Check(ctx context.Context, tx pgx.Tx, ownerID string, ids []string) error {
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		f, err := a.files.ByID(ctx, tx, id, false)
		if err != nil || f.OwnerID != ownerID || f.Status != files.StatusReady || seen[id] {
			return ErrBadFile
		}
		seen[id] = true
	}
	return nil
}

func (a Attachments) Load(ctx context.Context, tx pgx.Tx, ids []string) ([]files.File, error) {
	list := make([]files.File, 0, len(ids))
	for _, id := range ids {
		f, err := a.files.ByID(ctx, tx, id, false)
		if err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, nil
}
