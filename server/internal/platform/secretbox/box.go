package secretbox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

const KeySize = 32

var ErrBadKey = errors.New("secret key must be 32 bytes, base64 encoded")

var ErrCorrupt = errors.New("sealed secret is corrupt")

type Box struct {
	aead cipher.AEAD
}

func New(encodedKey string) (Box, error) {
	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil || len(key) != KeySize {
		return Box{}, ErrBadKey
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return Box{}, fmt.Errorf("cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return Box{}, fmt.Errorf("gcm: %w", err)
	}
	return Box{aead: aead}, nil
}

func (b Box) Seal(plain []byte) ([]byte, error) {
	nonce := make([]byte, b.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("nonce: %w", err)
	}
	return b.aead.Seal(nonce, nonce, plain, nil), nil
}

func (b Box) Open(sealed []byte) ([]byte, error) {
	size := b.aead.NonceSize()
	if len(sealed) < size {
		return nil, ErrCorrupt
	}
	plain, err := b.aead.Open(nil, sealed[:size], sealed[size:], nil)
	if err != nil {
		return nil, ErrCorrupt
	}
	return plain, nil
}
