package lti

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s LaunchService) DeepLinkForm(ctx context.Context, tenantSlug, code string) (string, string, error) {
	t, err := s.tenants.BySlug(ctx, tenantSlug)
	if err != nil {
		return "", "", err
	}
	var returnURL, response string
	err = s.tx.Run(ctx, t.ID, func(tx pgx.Tx) error {
		var err error
		returnURL, response, err = s.deepLinks.Take(ctx, tx, code)
		return err
	})
	return returnURL, response, err
}
