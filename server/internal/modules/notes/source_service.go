package notes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Sources(ctx context.Context, p auth.Principal) ([]Source, error) {
	var list []Source
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		rows, err := s.sources.Owned(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		list = make([]Source, 0, len(rows))
		for _, row := range rows {
			src, err := s.source(row)
			if err != nil {
				return err
			}
			list = append(list, src)
		}
		return nil
	})
	return list, err
}

func (s Service) SaveSource(ctx context.Context, p auth.Principal, in SourceInput) (Source, bool, error) {
	id := in.ID
	if !validID(id) {
		return Source{}, false, ErrInvalidSource
	}
	in, err := in.normalized()
	if err != nil {
		return Source{}, false, err
	}
	var src Source
	var created bool
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		existing, err := s.sources.Get(ctx, tx, id)
		switch {
		case err == nil && existing.OwnerID != p.UserID:
			return ErrSourceNotFound
		case err == nil:
		case errors.Is(err, ErrSourceNotFound):
			created = true
		default:
			return err
		}
		if in.FileID != nil {
			owner, err := s.sources.FileOwner(ctx, tx, *in.FileID)
			if err != nil {
				return err
			}
			if owner != p.UserID {
				return ErrInvalidSource
			}
		}
		raw, err := json.Marshal(in.SourceFields)
		if err != nil {
			return fmt.Errorf("encode source: %w", err)
		}
		data, err := s.sealer.seal(raw)
		if err != nil {
			return err
		}
		if err := s.checkQuota(ctx, tx, p.UserID, int64(len(data)-len(existing.Data))); err != nil {
			return err
		}
		row, err := s.sources.Upsert(ctx, tx, id, p.UserID, data, in.FileID)
		if err != nil {
			return err
		}
		src, err = s.source(row)
		return err
	})
	return src, created, err
}

func (s Service) DeleteSource(ctx context.Context, p auth.Principal, id string) error {
	if !validID(id) {
		return ErrSourceNotFound
	}
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		return s.sources.Delete(ctx, tx, id, p.UserID)
	})
}

func (s Service) LookupSource(ctx context.Context, in LookupInput) (SourceFields, error) {
	if (in.DOI == "") == (in.URL == "") {
		return SourceFields{}, ErrInvalidSource
	}
	if in.DOI != "" {
		return s.lookup.DOI(ctx, in.DOI)
	}
	return s.lookup.Page(ctx, in.URL)
}

func (s Service) source(row sourceRow) (Source, error) {
	raw, err := s.sealer.open(row.Data)
	if err != nil {
		return Source{}, err
	}
	var fields SourceFields
	if err := json.Unmarshal(raw, &fields); err != nil {
		return Source{}, fmt.Errorf("decode source: %w", err)
	}
	if fields.Authors == nil {
		fields.Authors = []string{}
	}
	return Source{ID: row.ID, SourceFields: fields, FileID: row.FileID, CreatedAt: row.CreatedAt, ChangedAt: row.ChangedAt}, nil
}
