package team

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

func (s Service) Card(ctx context.Context, p auth.Principal) (Card, error) {
	var c Card
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		c, err = s.cards.Get(ctx, tx, p.UserID)
		return err
	})
	return c, err
}

func (s Service) SaveCard(ctx context.Context, p auth.Principal, in Card) (Card, error) {
	in, err := in.normalized()
	if err != nil {
		return Card{}, err
	}
	var c Card
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.cards.Put(ctx, tx, p.UserID, in); err != nil {
			return err
		}
		c, err = s.cards.Get(ctx, tx, p.UserID)
		return err
	})
	return c, err
}

func (s Service) Suggest(ctx context.Context, p auth.Principal, kind, id string) (Suggestion, error) {
	if s.model == nil {
		return Suggestion{}, ErrAIUnavailable
	}
	var people []candidate
	var without []identity.Person
	var minSize, maxSize int
	var must string
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		ts, role, err := s.scope(ctx, tx, p, kind, id)
		if err != nil {
			return err
		}
		v, err := s.view(ctx, tx, p, ts, role)
		if err != nil {
			return err
		}
		minSize, maxSize = ts.MinSize, ts.MaxSize
		var pool []identity.Person
		switch {
		case role == RoleLead:
			pool = v.Unplaced
		case !ts.SelfFormed:
			return ErrNotSelfFormed
		case v.MyTeam != nil:
			return ErrInTeam
		case v.Solo:
			return ErrSeparated
		default:
			must = p.UserID
			me, err := s.people.ByIDs(ctx, tx, []string{p.UserID})
			if err != nil {
				return err
			}
			pool = me
			for _, seeker := range v.Seekers {
				if seeker.UserID != p.UserID {
					pool = append(pool, seeker.Person)
				}
			}
		}
		ids := make([]string, 0, len(pool))
		for _, person := range pool {
			ids = append(ids, person.UserID)
		}
		cards, err := s.cards.ByUsers(ctx, tx, ids)
		if err != nil {
			return err
		}
		for _, person := range pool {
			if card, ok := cards[person.UserID]; ok {
				people = append(people, candidate{person: person, card: card})
			} else {
				without = append(without, person)
			}
		}
		return nil
	})
	if err != nil {
		return Suggestion{}, err
	}
	out := Suggestion{Teams: []SuggestedTeam{}, Without: without}
	if without == nil {
		out.Without = []identity.Person{}
	}
	if len(people) < minSize || (must != "" && (len(people) == 0 || people[0].person.UserID != must)) {
		return out, nil
	}
	out.Teams, err = Matcher{model: s.model}.Match(ctx, people, minSize, maxSize, must)
	return out, err
}
