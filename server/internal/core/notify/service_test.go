package notify_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/live"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/notify"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type fixture struct {
	tx       database.TenantTx
	service  notify.Service
	notifier notify.Notifier
	me       auth.Principal
	other    auth.Principal
}

func setup(t *testing.T) fixture {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{tenant.Migrations(), identity.Migrations(), live.Migrations(), notify.Migrations()}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	users := identity.NewStore(tx)
	principal := func(subject string) auth.Principal {
		u, err := users.UpsertExternal(ctx, created.ID, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject})
		if err != nil {
			t.Fatalf("user: %v", err)
		}
		return auth.Principal{TenantID: created.ID, UserID: u.ID}
	}
	catalog, err := notify.NewCatalog(
		module.NotificationType{Key: "study.deadline", InApp: true, Push: true},
		module.NotificationType{Key: "core.security", Required: true},
	)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	return fixture{
		tx:       tx,
		service:  notify.NewService(tx, catalog),
		notifier: notify.NewNotifier(catalog),
		me:       principal("a"),
		other:    principal("b"),
	}
}

func (f fixture) send(t *testing.T, users ...string) {
	t.Helper()
	err := f.tx.Run(context.Background(), f.me.TenantID, func(tx pgx.Tx) error {
		return f.notifier.Send(context.Background(), tx, f.me.TenantID, notify.Message{Type: "study.deadline", UserIDs: users, Link: "/study"})
	})
	if err != nil {
		t.Fatalf("send: %v", err)
	}
}

func TestSendReachesOnlyRecipients(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	f.send(t, f.me.UserID)

	mine, err := f.service.List(ctx, f.me, nil, false)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	theirs, _ := f.service.List(ctx, f.other, nil, false)
	if len(mine) != 1 || mine[0].Type != "study.deadline" || len(theirs) != 0 {
		t.Fatalf("mine = %+v, theirs = %+v", mine, theirs)
	}
	var events []live.Event
	err = f.tx.Run(ctx, f.me.TenantID, func(tx pgx.Tx) error {
		var err error
		events, err = live.Store{}.After(ctx, tx, f.me.UserID, 0)
		return err
	})
	if err != nil || len(events) != 1 || events[0].Type != notify.LiveType {
		t.Fatalf("events = %+v, %v", events, err)
	}
}

func TestMarkReadUpdatesUnread(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	f.send(t, f.me.UserID)
	f.send(t, f.me.UserID)

	list, _ := f.service.List(ctx, f.me, nil, true)
	if err := f.service.MarkRead(ctx, f.me, list[0].ID); err != nil {
		t.Fatalf("mark: %v", err)
	}
	if n, _ := f.service.Unread(ctx, f.me); n != 1 {
		t.Fatalf("unread = %d, want 1", n)
	}
	if err := f.service.MarkAllRead(ctx, f.me); err != nil {
		t.Fatalf("mark all: %v", err)
	}
	if n, _ := f.service.Unread(ctx, f.me); n != 0 {
		t.Fatalf("unread = %d, want 0", n)
	}
}

func TestCannotReadSomeoneElsesNotification(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	f.send(t, f.me.UserID)
	list, _ := f.service.List(ctx, f.me, nil, false)

	if err := f.service.MarkRead(ctx, f.other, list[0].ID); !errors.Is(err, notify.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestTurnedOffTypeIsNotStored(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	if err := f.service.SetChannels(ctx, f.me, "study.deadline", notify.ChannelsInput{}); err != nil {
		t.Fatalf("set: %v", err)
	}
	f.send(t, f.me.UserID, f.other.UserID)

	if n, _ := f.service.Unread(ctx, f.me); n != 0 {
		t.Fatalf("unread = %d, want 0", n)
	}
	if n, _ := f.service.Unread(ctx, f.other); n != 1 {
		t.Fatalf("other unread = %d, want 1", n)
	}
}

func TestRequiredTypeStaysInApp(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	if err := f.service.SetChannels(ctx, f.me, "core.security", notify.ChannelsInput{}); !errors.Is(err, notify.ErrRequiredType) {
		t.Fatalf("err = %v, want ErrRequiredType", err)
	}
	channels, err := f.service.Channels(ctx, f.me)
	if err != nil {
		t.Fatalf("channels: %v", err)
	}
	if len(channels) != 2 || !channels[1].Required || !channels[1].InApp {
		t.Fatalf("channels = %+v", channels)
	}
}

func TestUnknownTypeIsRejected(t *testing.T) {
	f := setup(t)
	err := f.tx.Run(context.Background(), f.me.TenantID, func(tx pgx.Tx) error {
		return f.notifier.Send(context.Background(), tx, f.me.TenantID, notify.Message{Type: "nope.nope", UserIDs: []string{f.me.UserID}})
	})
	if !errors.Is(err, notify.ErrUnknownType) {
		t.Fatalf("err = %v, want ErrUnknownType", err)
	}
}
