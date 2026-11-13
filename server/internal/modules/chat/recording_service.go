package chat

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/livekit"
)

const recordingType = "video/mp4"

type RecordingService struct {
	k           module.Kernel
	chats       ChatStore
	calls       CallService
	store       CallStore
	recordings  RecordingStore
	files       files.FileStore
	messages    MessageStore
	attachments AttachmentStore
	announcer   Announcer
	now         func() time.Time
}

func NewRecordingService(k module.Kernel) RecordingService {
	sealer := Sealer{box: k.Secrets}
	return RecordingService{
		k: k, chats: ChatStore{sealer: sealer}, calls: NewCallService(k), messages: MessageStore{sealer: sealer}, now: time.Now,
	}
}

func (s RecordingService) Start(ctx context.Context, p auth.Principal, callID string) (Call, error) {
	if s.k.Calls == nil || !s.k.Calls.Records() {
		return Call{}, ErrRecordingUnavailable
	}
	if !validID(callID) {
		return Call{}, ErrCallNotFound
	}
	var fileID string
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		call, c, err := s.member(ctx, tx, p, callID)
		if err != nil {
			return err
		}
		if err := s.startable(ctx, tx, call, c); err != nil {
			return err
		}
		f, err := s.files.Create(ctx, tx, p.UserID, files.FileInput{
			Name: "call-" + s.now().UTC().Format("2006-01-02-1504") + ".mp4", ContentType: recordingType, Size: 1,
		})
		fileID = f.ID
		return err
	})
	if err != nil {
		return Call{}, err
	}
	egressID, err := s.k.Calls.StartRecording(ctx, RoomName(p.TenantID, callID), files.ObjectKey(p.TenantID, fileID), s.now())
	if err != nil {
		return Call{}, errors.Join(ErrRecordingUnavailable, err)
	}
	out, err := s.change(ctx, p, callID, func(tx pgx.Tx, call Call, c Chat) error {
		if err := s.startable(ctx, tx, call, c); err != nil {
			return err
		}
		return s.recordings.Start(ctx, tx, Recording{EgressID: egressID, CallID: call.ID, FileID: &fileID, StartedBy: &p.UserID})
	})
	if err != nil {
		return Call{}, errors.Join(err, s.k.Calls.StopRecording(context.WithoutCancel(ctx), egressID, s.now()))
	}
	return out, nil
}

func (s RecordingService) startable(ctx context.Context, tx pgx.Tx, call Call, c Chat) error {
	if call.Recording {
		return ErrRecordingState
	}
	allowed, err := s.calls.canRecord(ctx, tx, c)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrRecordingOff
	}
	return nil
}

func (s RecordingService) Stop(ctx context.Context, p auth.Principal, callID string) (Call, error) {
	if s.k.Calls == nil {
		return Call{}, ErrRecordingUnavailable
	}
	if !validID(callID) {
		return Call{}, ErrCallNotFound
	}
	var egressID string
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		call, c, err := s.member(ctx, tx, p, callID)
		if err != nil {
			return err
		}
		r, found, err := s.recordings.Active(ctx, tx, call.ID)
		if err != nil {
			return err
		}
		if !found {
			return ErrRecordingState
		}
		if (r.StartedBy == nil || *r.StartedBy != p.UserID) && !manages(c.Role) {
			return ErrRecordingOff
		}
		egressID = r.EgressID
		return nil
	})
	if err != nil {
		return Call{}, err
	}
	if err := s.k.Calls.StopRecording(ctx, egressID, s.now()); err != nil {
		return Call{}, errors.Join(ErrRecordingUnavailable, err)
	}
	return s.change(ctx, p, callID, func(tx pgx.Tx, _ Call, _ Chat) error {
		return s.recordings.End(ctx, tx, egressID)
	})
}

func (s RecordingService) Finish(ctx context.Context, tenantID string, info livekit.EgressInfo) error {
	return s.k.Tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		r, found, err := s.recordingByEgress(ctx, tx, info.EgressID)
		if err != nil || !found {
			return err
		}
		if err := s.recordings.End(ctx, tx, r.EgressID); err != nil {
			return err
		}
		_, active, err := s.store.ByID(ctx, tx, r.CallID)
		if err != nil {
			return err
		}
		if err := s.calls.announce(ctx, tx, tenantID, r.ChatID, r.CallID, active); err != nil {
			return err
		}
		if info.Status != livekit.EgressComplete || len(info.FileResults) == 0 || info.FileResults[0].Size <= 0 || r.FileID == nil || r.StartedBy == nil {
			return nil
		}
		if _, err := s.files.Finish(ctx, tx, *r.FileID, int64(info.FileResults[0].Size)); err != nil {
			if errors.Is(err, files.ErrNotFound) {
				return nil
			}
			return err
		}
		return s.post(ctx, tx, tenantID, r)
	})
}

func (s RecordingService) recordingByEgress(ctx context.Context, tx pgx.Tx, egressID string) (Recording, bool, error) {
	if egressID == "" {
		return Recording{}, false, nil
	}
	return s.recordings.ByEgress(ctx, tx, egressID)
}

func (s RecordingService) post(ctx context.Context, tx pgx.Tx, tenantID string, r Recording) error {
	seq, rev, err := s.chats.Next(ctx, tx, r.ChatID, true)
	if err != nil {
		return err
	}
	in := MessageInput{ClientID: "call-recording-" + r.EgressID, FileIDs: []string{*r.FileID}}
	m, err := s.messages.Insert(ctx, tx, r.ChatID, seq, rev, *r.StartedBy, in, messageBody{})
	if err != nil {
		return err
	}
	if err := s.attachments.Attach(ctx, tx, m.ID, in.FileIDs); err != nil {
		return err
	}
	list := []Message{m}
	if err := s.attachments.Fill(ctx, tx, list); err != nil {
		return err
	}
	return s.announcer.Message(ctx, tx, tenantID, list[0])
}

func (s RecordingService) change(ctx context.Context, p auth.Principal, callID string, fn func(pgx.Tx, Call, Chat) error) (Call, error) {
	if !validID(callID) {
		return Call{}, ErrCallNotFound
	}
	var out Call
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		call, c, err := s.member(ctx, tx, p, callID)
		if err != nil {
			return err
		}
		if err := fn(tx, call, c); err != nil {
			return err
		}
		if out, _, err = s.store.ByID(ctx, tx, callID); err != nil {
			return err
		}
		return s.calls.announce(ctx, tx, p.TenantID, call.ChatID, call.ID, true)
	})
	return out, err
}

func (s RecordingService) member(ctx context.Context, tx pgx.Tx, p auth.Principal, callID string) (Call, Chat, error) {
	call, found, err := s.store.ByID(ctx, tx, callID)
	if err != nil {
		return Call{}, Chat{}, err
	}
	if !found || !slices.Contains(call.People, p.UserID) {
		return Call{}, Chat{}, ErrCallNotFound
	}
	c, err := s.chats.ForMember(ctx, tx, p.UserID, call.ChatID)
	if err != nil {
		return Call{}, Chat{}, ErrCallNotFound
	}
	return call, c, nil
}
