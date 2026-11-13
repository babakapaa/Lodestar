package files

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

const (
	uploadTTL   = 15 * time.Minute
	downloadTTL = 5 * time.Minute
	megabyte    = 1 << 20
)

type Service struct {
	tenants  database.TenantTx
	storage  Storage
	settings Settings
	store    FileStore
	links    []Link
}

func NewService(tenants database.TenantTx, storage Storage, settings Settings, links ...Link) Service {
	return Service{tenants: tenants, storage: storage, settings: settings, links: links}
}

func (s Service) Create(ctx context.Context, p auth.Principal, in FileInput) (Created, error) {
	if s.storage == nil {
		return Created{}, ErrUnavailable
	}
	in, err := in.normalized()
	if err != nil {
		return Created{}, err
	}
	var created Created
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		limit, err := s.maxSize(ctx, tx)
		if err != nil {
			return err
		}
		if in.Size > limit {
			return ErrTooLarge
		}
		f, err := s.store.Create(ctx, tx, p.UserID, in)
		if err != nil {
			return err
		}
		upload, err := s.storage.PresignPut(ctx, ObjectKey(p.TenantID, f.ID), f.ContentType, f.Size, uploadTTL)
		if err != nil {
			return err
		}
		created = Created{File: f, Upload: upload}
		return nil
	})
	return created, err
}

func (s Service) Complete(ctx context.Context, p auth.Principal, id string) (File, error) {
	if s.storage == nil {
		return File{}, ErrUnavailable
	}
	var f File
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		if f, err = s.owned(ctx, tx, p, id, true); err != nil || f.Status == StatusReady {
			return err
		}
		size, err := s.storage.Size(ctx, ObjectKey(p.TenantID, f.ID))
		if errors.Is(err, objectstore.ErrMissing) {
			return ErrNotUploaded
		}
		if err != nil {
			return err
		}
		if size != f.Size {
			return ErrSizeMismatch
		}
		f, err = s.store.MarkReady(ctx, tx, f.ID)
		return err
	})
	return f, err
}

func (s Service) Get(ctx context.Context, p auth.Principal, id string) (File, error) {
	var f File
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		f, err = s.readable(ctx, tx, p, id)
		return err
	})
	return f, err
}

func (s Service) Download(ctx context.Context, p auth.Principal, id string) (Download, error) {
	if s.storage == nil {
		return Download{}, ErrUnavailable
	}
	f, err := s.Get(ctx, p, id)
	if err != nil {
		return Download{}, err
	}
	if f.Status != StatusReady {
		return Download{}, ErrNotUploaded
	}
	link, err := s.storage.PresignGet(ctx, ObjectKey(p.TenantID, f.ID), f.Name, f.ContentType, downloadTTL)
	if err != nil {
		return Download{}, err
	}
	return Download{URL: link}, nil
}

func (s Service) List(ctx context.Context, p auth.Principal, before *time.Time) ([]File, error) {
	var list []File
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.store.Page(ctx, tx, p.UserID, before)
		return err
	})
	return list, err
}

func (s Service) Delete(ctx context.Context, p auth.Principal, id string) error {
	if s.storage == nil {
		return ErrUnavailable
	}
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		f, err := s.owned(ctx, tx, p, id, true)
		if err != nil {
			return err
		}
		for _, link := range s.links {
			used, err := link.InUse(ctx, tx, f.ID)
			if err != nil {
				return err
			}
			if used {
				return ErrInUse
			}
		}
		if err := s.store.Delete(ctx, tx, f.ID); err != nil {
			return err
		}
		return s.storage.Delete(ctx, ObjectKey(p.TenantID, f.ID))
	})
}

func (s Service) owned(ctx context.Context, tx pgx.Tx, p auth.Principal, id string, lock bool) (File, error) {
	f, err := s.store.ByID(ctx, tx, id, lock)
	if err != nil {
		return File{}, err
	}
	if f.OwnerID != p.UserID {
		return File{}, ErrNotFound
	}
	return f, nil
}

func (s Service) readable(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (File, error) {
	f, err := s.store.ByID(ctx, tx, id, false)
	if err != nil {
		return File{}, err
	}
	if f.OwnerID == p.UserID {
		return f, nil
	}
	for _, link := range s.links {
		ok, err := link.Readable(ctx, tx, p, f.ID)
		if err != nil {
			return File{}, err
		}
		if ok {
			return f, nil
		}
	}
	return File{}, ErrNotFound
}

func (s Service) maxSize(ctx context.Context, tx pgx.Tx) (int64, error) {
	raw, err := s.settings.Value(ctx, tx, MaxFileSetting)
	if err != nil {
		return 0, err
	}
	var mb int64
	if err := json.Unmarshal(raw, &mb); err != nil {
		return 0, fmt.Errorf("read %s: %w", MaxFileSetting, err)
	}
	return mb * megabyte, nil
}
