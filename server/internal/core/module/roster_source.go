package module

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrRosterNotFound = errors.New("roster not found")

type RosterSource interface {
	Kind() string
	All(ctx context.Context, tx pgx.Tx) ([]Roster, error)
	One(ctx context.Context, tx pgx.Tx, id string) (Roster, error)
}
