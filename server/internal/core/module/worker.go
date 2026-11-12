package module

import "context"

type Worker interface {
	Run(ctx context.Context)
}
