package files

import (
	"context"
	"io"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/privacy"
)

type PrivacySection struct {
	storage Storage
	store   FileStore
}

func NewPrivacySection(storage Storage) PrivacySection {
	return PrivacySection{storage: storage}
}

func (PrivacySection) Name() string {
	return "files"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	list, err := s.store.Owned(ctx, tx, userID)
	if list == nil {
		list = []File{}
	}
	return list, err
}

func (s PrivacySection) Blobs(ctx context.Context, tx pgx.Tx, tenantID, userID string) ([]privacy.Blob, error) {
	if s.storage == nil {
		return nil, nil
	}
	list, err := s.store.Owned(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	var blobs []privacy.Blob
	for _, f := range list {
		if f.Status != StatusReady {
			continue
		}
		key := ObjectKey(tenantID, f.ID)
		blobs = append(blobs, privacy.Blob{
			Path: "files/" + f.ID + "/" + f.Name,
			Open: func(ctx context.Context) (io.ReadCloser, error) { return s.storage.Open(ctx, key) },
		})
	}
	return blobs, nil
}

func (s PrivacySection) Erase(ctx context.Context, tx pgx.Tx, tenantID, userID string) error {
	list, err := s.store.Owned(ctx, tx, userID)
	if err != nil {
		return err
	}
	if len(list) > 0 && s.storage == nil {
		return ErrUnavailable
	}
	for _, f := range list {
		if err := s.storage.Delete(ctx, ObjectKey(tenantID, f.ID)); err != nil {
			return err
		}
	}
	return s.store.DeleteOwned(ctx, tx, userID)
}
