package app

import "context"

type Worker interface {
	Run(ctx context.Context)
}
