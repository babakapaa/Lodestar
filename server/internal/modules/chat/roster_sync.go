package chat

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type RosterSync struct {
	chats     ChatStore
	members   MemberStore
	announcer Announcer
}

func (s RosterSync) Sync(ctx context.Context, tx pgx.Tx, tenantID string, r module.Roster) error {
	origin := Origin{Kind: r.Kind, ID: r.ID}
	if !r.Archived {
		if _, _, err := s.chats.Insert(ctx, tx, KindChannel, r.Title, "", &origin, nil, nil); err != nil {
			return err
		}
	}
	refs, err := s.chats.ByOrigin(ctx, tx, origin)
	if err != nil {
		return err
	}
	for _, ref := range refs {
		if ref.Archived != r.Archived {
			if err := s.chats.Archive(ctx, tx, ref.ID, r.Archived); err != nil {
				return err
			}
			if err := s.announcer.Changed(ctx, tx, tenantID, ref.ID); err != nil {
				return err
			}
		}
		if r.Archived {
			continue
		}
		changed, err := s.members.Sync(ctx, tx, ref.ID, r.Leads, r.Members)
		if err != nil {
			return err
		}
		if err := s.announcer.ChangedFor(ctx, tx, tenantID, ref.ID, changed); err != nil {
			return err
		}
	}
	return nil
}
