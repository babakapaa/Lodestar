package lti

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

const keyBits = 2048

type KeyStore struct {
	box secretbox.Box
}

func NewKeyStore(box secretbox.Box) KeyStore {
	return KeyStore{box: box}
}

func (s KeyStore) Current(ctx context.Context, tx pgx.Tx) (Key, error) {
	var k Key
	var sealed []byte
	err := tx.QueryRow(ctx, `SELECT kid, private_key FROM lti.keys ORDER BY created_at DESC, kid LIMIT 1`).Scan(&k.ID, &sealed)
	if errors.Is(err, pgx.ErrNoRows) {
		return s.create(ctx, tx)
	}
	if err != nil {
		return Key{}, fmt.Errorf("find lti key: %w", err)
	}
	return s.open(k.ID, sealed)
}

func (s KeyStore) All(ctx context.Context, tx pgx.Tx) ([]Key, error) {
	if _, err := s.Current(ctx, tx); err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, `SELECT kid, private_key FROM lti.keys ORDER BY created_at DESC, kid`)
	if err != nil {
		return nil, fmt.Errorf("list lti keys: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Key, error) {
		var kid string
		var sealed []byte
		if err := row.Scan(&kid, &sealed); err != nil {
			return Key{}, err
		}
		return s.open(kid, sealed)
	})
}

func (s KeyStore) create(ctx context.Context, tx pgx.Tx) (Key, error) {
	private, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return Key{}, fmt.Errorf("generate lti key: %w", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return Key{}, fmt.Errorf("encode lti key: %w", err)
	}
	sealed, err := s.box.Seal(der)
	if err != nil {
		return Key{}, err
	}
	kid, err := randomString()
	if err != nil {
		return Key{}, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO lti.keys (kid, private_key) VALUES ($1, $2)`, kid, sealed); err != nil {
		return Key{}, fmt.Errorf("save lti key: %w", err)
	}
	return Key{ID: kid, Private: private}, nil
}

func (s KeyStore) open(kid string, sealed []byte) (Key, error) {
	der, err := s.box.Open(sealed)
	if err != nil {
		return Key{}, fmt.Errorf("open lti key: %w", err)
	}
	parsed, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return Key{}, fmt.Errorf("parse lti key: %w", err)
	}
	private, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return Key{}, fmt.Errorf("lti key %s is not rsa", kid)
	}
	return Key{ID: kid, Private: private}, nil
}
