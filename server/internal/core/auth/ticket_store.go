package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type TicketStore struct {
	tenants database.TenantTx
	ttl     time.Duration
}

func NewTicketStore(tenants database.TenantTx, ttl time.Duration) TicketStore {
	return TicketStore{tenants: tenants, ttl: ttl}
}

func (s TicketStore) Create(ctx context.Context, tenantID, userID, device string) (string, error) {
	ticket, err := newToken(ticketPrefix, tenantID)
	if err != nil {
		return "", err
	}
	err = s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `DELETE FROM core.launch_tickets WHERE expires_at < now()`); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO core.launch_tickets (ticket_hash, user_id, device_name, expires_at)
			VALUES ($1, $2, $3, now() + $4::interval)`, ticket.Hash, userID, deviceName(device), s.ttl)
		return err
	})
	if err != nil {
		return "", fmt.Errorf("save launch ticket: %w", err)
	}
	return ticket.Value, nil
}

func (s TicketStore) Take(ctx context.Context, value string) (LoginGrant, error) {
	ticket, err := parseToken(ticketPrefix, value)
	if err != nil {
		return LoginGrant{}, err
	}
	g := LoginGrant{TenantID: ticket.TenantID}
	err = s.tenants.Run(ctx, ticket.TenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `
			DELETE FROM core.launch_tickets WHERE ticket_hash = $1 AND expires_at > now()
			RETURNING user_id::text, device_name`, ticket.Hash).Scan(&g.UserID, &g.DeviceName)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return LoginGrant{}, ErrInvalidToken
	}
	if err != nil {
		return LoginGrant{}, fmt.Errorf("take launch ticket: %w", err)
	}
	return g, nil
}
