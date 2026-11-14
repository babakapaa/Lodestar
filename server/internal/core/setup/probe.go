package setup

import "context"

type Probe interface {
	Name() string
	Required() bool
	Check(ctx context.Context) error
}
