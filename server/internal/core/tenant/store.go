package tenant

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueViolation = "23505"
	checkViolation  = "23514"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return Store{pool: pool}
}

func (s Store) Create(ctx context.Context, slug, name string) (Tenant, error) {
	row := s.pool.QueryRow(ctx,
		`INSERT INTO core.tenants (slug, name) VALUES ($1, $2)
		 RETURNING id::text, slug, name, created_at`, slug, name)
	t, err := scan(row)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case uniqueViolation:
			return Tenant{}, ErrSlugTaken
		case checkViolation:
			return Tenant{}, ErrInvalid
		}
	}
	if err != nil {
		return Tenant{}, fmt.Errorf("create tenant: %w", err)
	}
	return t, nil
}

func (s Store) BySlug(ctx context.Context, slug string) (Tenant, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id::text, slug, name, created_at FROM core.tenants WHERE slug = $1`, slug)
	t, err := scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Tenant{}, ErrNotFound
	}
	if err != nil {
		return Tenant{}, fmt.Errorf("find tenant: %w", err)
	}
	return t, nil
}

func (s Store) ByID(ctx context.Context, id string) (Tenant, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id::text, slug, name, created_at FROM core.tenants WHERE id::text = $1`, id)
	t, err := scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Tenant{}, ErrNotFound
	}
	if err != nil {
		return Tenant{}, fmt.Errorf("find tenant: %w", err)
	}
	return t, nil
}

func scan(row pgx.Row) (Tenant, error) {
	var t Tenant
	err := row.Scan(&t.ID, &t.Slug, &t.Name, &t.CreatedAt)
	return t, err
}
