package lti

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const deepLinkLife = time.Hour

type DeepLinkStore struct{}

func (DeepLinkStore) Create(ctx context.Context, tx pgx.Tx, d DeepLink) (string, error) {
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO lti.deep_links (platform_id, user_id, deployment_id, return_url, data, accept_multiple)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		d.PlatformID, d.UserID, d.DeploymentID, d.ReturnURL, d.Data, d.AcceptMultiple).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("create deep link: %w", err)
	}
	return id, nil
}

func (DeepLinkStore) Open(ctx context.Context, tx pgx.Tx, id, userID string) (DeepLink, error) {
	var d DeepLink
	err := tx.QueryRow(ctx, `
		SELECT d.id, d.platform_id, d.user_id, d.deployment_id, d.return_url, d.data, d.accept_multiple, d.created_at, p.name
		FROM lti.deep_links d JOIN lti.platforms p ON p.id = d.platform_id
		WHERE d.id::text = $1 AND d.user_id = $2 AND d.response IS NULL AND d.created_at > $3`,
		id, userID, time.Now().Add(-deepLinkLife)).
		Scan(&d.ID, &d.PlatformID, &d.UserID, &d.DeploymentID, &d.ReturnURL, &d.Data, &d.AcceptMultiple, &d.CreatedAt, &d.Platform)
	if errors.Is(err, pgx.ErrNoRows) {
		return DeepLink{}, ErrDeepLinkNotFound
	}
	if err != nil {
		return DeepLink{}, fmt.Errorf("read deep link: %w", err)
	}
	return d, nil
}

func (DeepLinkStore) Answer(ctx context.Context, tx pgx.Tx, id, response, code string) error {
	if _, err := tx.Exec(ctx, `UPDATE lti.deep_links SET response = $2, code = $3 WHERE id = $1`, id, response, code); err != nil {
		return fmt.Errorf("answer deep link: %w", err)
	}
	return nil
}

func (DeepLinkStore) Take(ctx context.Context, tx pgx.Tx, code string) (string, string, error) {
	var returnURL, response string
	err := tx.QueryRow(ctx, `
		DELETE FROM lti.deep_links WHERE code = $1 AND response IS NOT NULL AND created_at > $2
		RETURNING return_url, response`, code, time.Now().Add(-deepLinkLife)).Scan(&returnURL, &response)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", ErrDeepLinkNotFound
	}
	if err != nil {
		return "", "", fmt.Errorf("take deep link: %w", err)
	}
	return returnURL, response, nil
}
