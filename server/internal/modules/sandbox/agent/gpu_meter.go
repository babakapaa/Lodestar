package agent

import "context"

type GPUMeter interface {
	Load(ctx context.Context) (map[int]float64, error)
}
