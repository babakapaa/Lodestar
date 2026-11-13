package module

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TeamScopes []TeamScopeSource

func (s TeamScopes) Find(kind string) (TeamScopeSource, bool) {
	for _, source := range s {
		if source.Kind() == kind {
			return source, true
		}
	}
	return nil, false
}

func (s TeamScopes) Load(ctx context.Context, tx pgx.Tx, kind, id string) (TeamScope, error) {
	source, ok := s.Find(kind)
	if !ok {
		return TeamScope{}, ErrTeamScopeNotFound
	}
	return source.TeamScope(ctx, tx, id)
}
