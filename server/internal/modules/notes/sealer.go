package notes

import (
	"encoding/json"
	"fmt"

	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type Sealer struct {
	box secretbox.Box
}

func (s Sealer) seal(plain []byte) ([]byte, error) {
	sealed, err := s.box.Seal(plain)
	if err != nil {
		return nil, fmt.Errorf("seal note data: %w", err)
	}
	return sealed, nil
}

func (s Sealer) open(sealed []byte) ([]byte, error) {
	plain, err := s.box.Open(sealed)
	if err != nil {
		return nil, fmt.Errorf("open note data: %w", err)
	}
	return plain, nil
}

func (s Sealer) sealMeta(m Meta) ([]byte, error) {
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return s.seal(raw)
}

func (s Sealer) openMeta(sealed []byte) (Meta, error) {
	raw, err := s.open(sealed)
	if err != nil {
		return Meta{}, err
	}
	var m Meta
	if err := json.Unmarshal(raw, &m); err != nil {
		return Meta{}, fmt.Errorf("decode note meta: %w", err)
	}
	return m, nil
}
