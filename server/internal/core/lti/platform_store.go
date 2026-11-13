package lti

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolation = "23505"
	checkViolation  = "23514"
	invalidText     = "22P02"
)

const platformColumns = `id::text, name, issuer, client_id, auth_url, token_url, jwks_url, deployment_ids, match_email, created_at`

type PlatformStore struct{}

func (PlatformStore) List(ctx context.Context, tx pgx.Tx) ([]Platform, error) {
	rows, err := tx.Query(ctx, `SELECT `+platformColumns+` FROM lti.platforms ORDER BY name, id`)
	if err != nil {
		return nil, fmt.Errorf("list lti platforms: %w", err)
	}
	return pgx.CollectRows(rows, scanPlatform)
}

func (PlatformStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Platform, error) {
	rows, err := tx.Query(ctx, `SELECT `+platformColumns+` FROM lti.platforms WHERE id::text = $1`, id)
	if err != nil {
		return Platform{}, fmt.Errorf("find lti platform: %w", err)
	}
	p, err := pgx.CollectExactlyOneRow(rows, scanPlatform)
	if errors.Is(err, pgx.ErrNoRows) {
		return Platform{}, ErrPlatformNotFound
	}
	return p, err
}

func (PlatformStore) ByIssuer(ctx context.Context, tx pgx.Tx, issuer, clientID string) (Platform, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+platformColumns+` FROM lti.platforms
		WHERE issuer = $1 AND ($2 = '' OR client_id = $2)`, issuer, clientID)
	if err != nil {
		return Platform{}, fmt.Errorf("find lti platform: %w", err)
	}
	list, err := pgx.CollectRows(rows, scanPlatform)
	if err != nil {
		return Platform{}, fmt.Errorf("find lti platform: %w", err)
	}
	switch len(list) {
	case 0:
		return Platform{}, ErrPlatformNotFound
	case 1:
		return list[0], nil
	default:
		return Platform{}, ErrPlatformAmbiguous
	}
}

func (PlatformStore) Insert(ctx context.Context, tx pgx.Tx, in PlatformInput) (Platform, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO lti.platforms (name, issuer, client_id, auth_url, token_url, jwks_url, deployment_ids, match_email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING `+platformColumns,
		in.Name, in.Issuer, in.ClientID, in.AuthURL, in.TokenURL, in.JWKSURL, in.DeploymentIDs, in.MatchEmail)
	if err != nil {
		return Platform{}, platformError(err)
	}
	p, err := pgx.CollectExactlyOneRow(rows, scanPlatform)
	return p, platformError(err)
}

func (PlatformStore) Update(ctx context.Context, tx pgx.Tx, id string, in PlatformInput) (Platform, error) {
	rows, err := tx.Query(ctx, `
		UPDATE lti.platforms
		SET name = $2, issuer = $3, client_id = $4, auth_url = $5, token_url = $6, jwks_url = $7,
		    deployment_ids = $8, match_email = $9, updated_at = now()
		WHERE id::text = $1
		RETURNING `+platformColumns,
		id, in.Name, in.Issuer, in.ClientID, in.AuthURL, in.TokenURL, in.JWKSURL, in.DeploymentIDs, in.MatchEmail)
	if err != nil {
		return Platform{}, platformError(err)
	}
	p, err := pgx.CollectExactlyOneRow(rows, scanPlatform)
	if errors.Is(err, pgx.ErrNoRows) {
		return Platform{}, ErrPlatformNotFound
	}
	return p, platformError(err)
}

func (PlatformStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	tag, err := tx.Exec(ctx, `DELETE FROM lti.platforms WHERE id::text = $1`, id)
	if err != nil {
		return fmt.Errorf("delete lti platform: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrPlatformNotFound
	}
	return nil
}

func scanPlatform(row pgx.CollectableRow) (Platform, error) {
	var p Platform
	err := row.Scan(&p.ID, &p.Name, &p.Issuer, &p.ClientID, &p.AuthURL, &p.TokenURL, &p.JWKSURL, &p.DeploymentIDs, &p.MatchEmail, &p.CreatedAt)
	return p, err
}

func platformError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case uniqueViolation:
			return ErrPlatformTaken
		case checkViolation, invalidText:
			return ErrInvalidPlatform
		}
	}
	if err != nil {
		return fmt.Errorf("save lti platform: %w", err)
	}
	return nil
}
