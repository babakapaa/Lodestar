package profile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ProfileStore struct{}

func (ProfileStore) Get(ctx context.Context, tx pgx.Tx, userID string) (Profile, error) {
	p := Profile{UserID: userID}
	var contacts []byte
	err := tx.QueryRow(ctx, `
		SELECT coalesce(nullif(p.display_name, ''), u.display_name), coalesce(p.about, ''), coalesce(p.contacts, '[]')
		FROM core.users u LEFT JOIN core.profiles p ON p.user_id = u.id
		WHERE u.id::text = $1`, userID).Scan(&p.DisplayName, &p.About, &contacts)
	if errors.Is(err, pgx.ErrNoRows) {
		return Profile{}, ErrUserNotFound
	}
	if err != nil {
		return Profile{}, fmt.Errorf("read profile: %w", err)
	}
	if err := json.Unmarshal(contacts, &p.Contacts); err != nil {
		return Profile{}, fmt.Errorf("decode contacts: %w", err)
	}
	return p, nil
}

func (ProfileStore) Put(ctx context.Context, tx pgx.Tx, userID string, in ProfileInput) error {
	contacts, err := json.Marshal(contactsOrEmpty(in.Contacts))
	if err != nil {
		return fmt.Errorf("encode contacts: %w", err)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO core.profiles (user_id, display_name, about, contacts) VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET display_name = EXCLUDED.display_name, about = EXCLUDED.about, contacts = EXCLUDED.contacts, updated_at = now()`,
		userID, in.DisplayName, in.About, contacts)
	if err != nil {
		return fmt.Errorf("save profile: %w", err)
	}
	return nil
}

func contactsOrEmpty(c []Contact) []Contact {
	if c == nil {
		return []Contact{}
	}
	return c
}
