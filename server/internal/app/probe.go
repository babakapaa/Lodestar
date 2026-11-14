package app

import "context"

type probe struct {
	name     string
	required bool
	check    func(ctx context.Context) error
}

func (p probe) Name() string {
	return p.name
}

func (p probe) Required() bool {
	return p.required
}

func (p probe) Check(ctx context.Context) error {
	return p.check(ctx)
}
