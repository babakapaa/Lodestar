package notes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/live"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const ChangedEvent = "notes.changed"

type Service struct {
	tenants  database.TenantTx
	settings module.Settings
	sealer   Sealer
	notes    NoteStore
	updates  UpdateStore
	versions VersionStore
	members  MemberStore
	people   identity.People
	live     live.Store
}

func NewService(k module.Kernel) Service {
	return Service{tenants: k.Tenants, settings: k.Settings, sealer: Sealer{box: k.Secrets}}
}

func (s Service) List(ctx context.Context, p auth.Principal) ([]Note, error) {
	var list []Note
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		rows, err := s.notes.Visible(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		list = make([]Note, 0, len(rows))
		for _, row := range rows {
			n, err := s.note(row)
			if err != nil {
				return err
			}
			list = append(list, n)
		}
		return nil
	})
	return list, err
}

func (s Service) Get(ctx context.Context, p auth.Principal, id string) (Note, error) {
	var n Note
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		row, err := s.notes.Get(ctx, tx, p.UserID, id, false)
		if err != nil {
			return err
		}
		n, err = s.note(row)
		return err
	})
	return n, err
}

func (s Service) Create(ctx context.Context, p auth.Principal, in CreateInput) (Note, bool, error) {
	in, err := in.normalized()
	if err != nil {
		return Note{}, false, err
	}
	var n Note
	var created bool
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		meta, err := s.sealer.sealMeta(in.Meta)
		if err != nil {
			return err
		}
		if err := s.checkQuota(ctx, tx, p.UserID, int64(len(meta))); err != nil {
			return err
		}
		if created, err = s.notes.Insert(ctx, tx, in.ID, p.UserID, meta); err != nil {
			return err
		}
		row, err := s.notes.Get(ctx, tx, p.UserID, in.ID, false)
		if errors.Is(err, ErrNoteNotFound) || (err == nil && row.Role != RoleOwner) {
			return ErrNoteExists
		}
		if err != nil {
			return err
		}
		if n, err = s.note(row); err != nil {
			return err
		}
		if created {
			return s.announce(ctx, tx, p, row.ID, row.Seq)
		}
		return nil
	})
	return n, created, err
}

func (s Service) SetMeta(ctx context.Context, p auth.Principal, id string, in Meta) (Note, error) {
	in, err := in.normalized()
	if err != nil {
		return Note{}, err
	}
	var n Note
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		row, err := s.editable(ctx, tx, p, id)
		if err != nil {
			return err
		}
		meta, err := s.sealer.sealMeta(in)
		if err != nil {
			return err
		}
		delta := int64(len(meta) - len(row.Meta))
		if err := s.checkQuota(ctx, tx, row.OwnerID, delta); err != nil {
			return err
		}
		if err := s.notes.SetMeta(ctx, tx, id, meta, delta); err != nil {
			return err
		}
		if row, err = s.notes.Get(ctx, tx, p.UserID, id, false); err != nil {
			return err
		}
		if n, err = s.note(row); err != nil {
			return err
		}
		return s.announce(ctx, tx, p, id, row.Seq)
	})
	return n, err
}

func (s Service) SetDeleted(ctx context.Context, p auth.Principal, id string, deleted bool) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		row, err := s.notes.Get(ctx, tx, p.UserID, id, true)
		if err != nil {
			return err
		}
		if row.Role != RoleOwner {
			return ErrNotOwner
		}
		if err := s.notes.SetDeleted(ctx, tx, id, deleted); err != nil {
			return err
		}
		return s.announce(ctx, tx, p, id, row.Seq)
	})
}

func (s Service) Changes(ctx context.Context, p auth.Principal, id string, after int64) (Changes, error) {
	var c Changes
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		row, err := s.notes.Get(ctx, tx, p.UserID, id, false)
		if err != nil {
			return err
		}
		c.Seq = row.Seq
		from := after
		if after < row.BaseSeq {
			latest, ok, err := s.versions.Latest(ctx, tx, id)
			if err != nil {
				return err
			}
			if ok {
				data, err := s.sealer.open(latest.Data)
				if err != nil {
					return err
				}
				c.Snapshot = &Snapshot{Seq: latest.Seq, Data: data}
				from = latest.Seq
			}
		}
		updates, err := s.updates.After(ctx, tx, id, from)
		if err != nil {
			return err
		}
		for i := range updates {
			if updates[i].Data, err = s.sealer.open(updates[i].Data); err != nil {
				return err
			}
		}
		c.Updates = updates
		return nil
	})
	return c, err
}

func (s Service) Push(ctx context.Context, p auth.Principal, id string, in UpdateInput) (Update, bool, error) {
	if !in.valid() {
		return Update{}, false, ErrInvalidUpdate
	}
	var u Update
	var created bool
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		row, err := s.editable(ctx, tx, p, id)
		if err != nil {
			return err
		}
		existing, found, err := s.updates.ByKey(ctx, tx, id, in.Key)
		if err != nil {
			return err
		}
		if found {
			existing.Data = in.Data
			u = existing
			return nil
		}
		sealed, err := s.sealer.seal(in.Data)
		if err != nil {
			return err
		}
		if err := s.checkQuota(ctx, tx, row.OwnerID, int64(len(sealed))); err != nil {
			return err
		}
		seq, err := s.notes.Advance(ctx, tx, id, int64(len(sealed)))
		if err != nil {
			return err
		}
		if err := s.updates.Insert(ctx, tx, id, seq, in.Key, p.UserID, sealed); err != nil {
			return err
		}
		author := p.UserID
		u = Update{Seq: seq, Key: in.Key, AuthorID: &author, Data: in.Data}
		created = true
		return s.announce(ctx, tx, p, id, seq)
	})
	return u, created, err
}

func (s Service) Compact(ctx context.Context, p auth.Principal, id string, in SnapshotInput) (Version, error) {
	if !in.valid() {
		return Version{}, ErrInvalidUpdate
	}
	var v Version
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		row, err := s.editable(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if in.Seq > row.Seq || in.Seq <= row.BaseSeq {
			return ErrStaleSnapshot
		}
		data, err := s.sealer.seal(in.Data)
		if err != nil {
			return err
		}
		body, err := s.sealer.seal([]byte(in.Body))
		if err != nil {
			return err
		}
		added := int64(len(data) + len(body))
		if err := s.checkQuota(ctx, tx, row.OwnerID, added); err != nil {
			return err
		}
		if v, err = s.versions.Insert(ctx, tx, id, in.Seq, p.UserID, data, body); err != nil {
			return err
		}
		freed, err := s.updates.Prune(ctx, tx, id, in.Seq)
		if err != nil {
			return err
		}
		return s.notes.Rebase(ctx, tx, id, in.Seq, added-freed)
	})
	return v, err
}

func (s Service) Versions(ctx context.Context, p auth.Principal, id string) ([]Version, error) {
	var list []Version
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.notes.Get(ctx, tx, p.UserID, id, false); err != nil {
			return err
		}
		var err error
		list, err = s.versions.List(ctx, tx, id)
		return err
	})
	return list, err
}

func (s Service) Version(ctx context.Context, p auth.Principal, id, versionID string) (VersionContent, error) {
	var c VersionContent
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.notes.Get(ctx, tx, p.UserID, id, false); err != nil {
			return err
		}
		row, err := s.versions.Get(ctx, tx, id, versionID)
		if err != nil {
			return err
		}
		data, err := s.sealer.open(row.Data)
		if err != nil {
			return err
		}
		body, err := s.sealer.open(row.Body)
		if err != nil {
			return err
		}
		c = VersionContent{Version: row.Version, Data: data, Body: string(body)}
		return nil
	})
	return c, err
}

func (s Service) Members(ctx context.Context, p auth.Principal, id string) ([]Member, error) {
	var list []Member
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.notes.Get(ctx, tx, p.UserID, id, false); err != nil {
			return err
		}
		var err error
		list, err = s.members.List(ctx, tx, id)
		return err
	})
	return list, err
}

func (s Service) SetMember(ctx context.Context, p auth.Principal, id, userID string, in MemberInput) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		row, err := s.notes.Get(ctx, tx, p.UserID, id, true)
		if err != nil {
			return err
		}
		if row.Role != RoleOwner {
			return ErrNotOwner
		}
		if userID == row.OwnerID {
			return ErrInvalidNote
		}
		exists, err := s.people.Exists(ctx, tx, userID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrUserNotFound
		}
		if err := s.members.Put(ctx, tx, id, userID, in.CanEdit, p.UserID); err != nil {
			return err
		}
		return s.announce(ctx, tx, p, id, row.Seq)
	})
}

func (s Service) RemoveMember(ctx context.Context, p auth.Principal, id, userID string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		row, err := s.notes.Get(ctx, tx, p.UserID, id, true)
		if err != nil {
			return err
		}
		if row.Role != RoleOwner && userID != p.UserID {
			return ErrNotOwner
		}
		if err := s.announce(ctx, tx, p, id, row.Seq); err != nil {
			return err
		}
		return s.members.Remove(ctx, tx, id, userID)
	})
}

func (s Service) editable(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (noteRow, error) {
	row, err := s.notes.Get(ctx, tx, p.UserID, id, true)
	if err != nil {
		return noteRow{}, err
	}
	if !row.Role.canEdit() {
		return noteRow{}, ErrReadOnly
	}
	if row.DeletedAt != nil {
		return noteRow{}, ErrNoteDeleted
	}
	return row, nil
}

func (s Service) checkQuota(ctx context.Context, tx pgx.Tx, ownerID string, added int64) error {
	if added <= 0 {
		return nil
	}
	raw, err := s.settings.Value(ctx, tx, QuotaSetting)
	if err != nil {
		return err
	}
	var mb int64
	if err := json.Unmarshal(raw, &mb); err != nil {
		return fmt.Errorf("decode notes quota: %w", err)
	}
	used, err := s.notes.OwnerBytes(ctx, tx, ownerID)
	if err != nil {
		return err
	}
	if used+added > mb<<20 {
		return ErrQuotaExceeded
	}
	return nil
}

func (s Service) announce(ctx context.Context, tx pgx.Tx, p auth.Principal, id string, seq int64) error {
	audience, err := s.notes.Audience(ctx, tx, id)
	if err != nil {
		return err
	}
	return s.live.Publish(ctx, tx, p.TenantID, audience, ChangedEvent, map[string]any{"noteId": id, "seq": seq})
}

func (s Service) note(row noteRow) (Note, error) {
	meta, err := s.sealer.openMeta(row.Meta)
	if err != nil {
		return Note{}, err
	}
	return Note{
		ID: row.ID, OwnerID: row.OwnerID, Role: row.Role, Meta: meta, Seq: row.Seq, Bytes: row.Bytes,
		CreatedAt: row.CreatedAt, ChangedAt: row.ChangedAt, DeletedAt: row.DeletedAt,
	}, nil
}
