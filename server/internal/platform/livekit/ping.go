package livekit

import (
	"context"
	"time"
)

func (i *Issuer) Ping(ctx context.Context) error {
	return i.call(ctx, "ListEgress", map[string]any{}, nil, time.Now())
}
