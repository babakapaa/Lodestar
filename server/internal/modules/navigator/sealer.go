package navigator

import (
	"encoding/json"
	"fmt"

	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type Sealer struct {
	box secretbox.Box
}

func (s Sealer) seal(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	sealed, err := s.box.Seal(raw)
	if err != nil {
		return nil, fmt.Errorf("seal navigator data: %w", err)
	}
	return sealed, nil
}

func (s Sealer) open(sealed []byte, v any) error {
	raw, err := s.box.Open(sealed)
	if err != nil {
		return fmt.Errorf("open navigator data: %w", err)
	}
	if err := json.Unmarshal(raw, v); err != nil {
		return fmt.Errorf("decode navigator data: %w", err)
	}
	return nil
}
