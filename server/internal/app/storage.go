package app

import (
	"context"
	"log/slog"

	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

func OpenStorage(ctx context.Context, cfg objectstore.Config) (files.Storage, error) {
	if !cfg.Enabled() {
		slog.Warn("file storage is not configured, set LODESTAR_S3_ENDPOINT to enable files")
		return nil, nil
	}
	store, err := objectstore.New(cfg)
	if err != nil {
		return nil, err
	}
	if err := store.EnsureBucket(ctx); err != nil {
		return nil, err
	}
	return store, nil
}
