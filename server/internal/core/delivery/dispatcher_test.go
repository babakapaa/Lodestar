package delivery_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/delivery"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/live"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/notify"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/fcm"
)

var deadline = module.NotificationType{
	Key:   "study.deadline",
	Push:  true,
	Email: true,
	Names: map[string]string{"en": "Deadlines"},
	Texts: map[string]module.NotificationText{
		"ru": {Title: "Срок сдачи", Body: "Осталось {days} дн."},
		"kk": {Title: "Тапсыру мерзімі", Body: "{days} күн қалды"},
		"en": {Title: "Deadline", Body: "{days} days left"},
	},
}

type fixture struct {
	pool     *pgxpool.Pool
	tx       database.TenantTx
	devices  delivery.Service
	me       auth.Principal
	other    auth.Principal
	pusher   *fakePusher
	mailer   *fakeMailer
	prefs    values
	settings values
}

func setup(t *testing.T) *fixture {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{tenant.Migrations(), identity.Migrations(), live.Migrations(), notify.Migrations(), delivery.Migrations()}
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
		u, err := users.UpsertExternal(ctx, created.ID, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject, Email: subject + "@uni.example"})
		if err != nil {
			t.Fatalf("user: %v", err)
		}
		return auth.Principal{TenantID: created.ID, UserID: u.ID}
	}
	return &fixture{
		pool:    pool,
		tx:      tx,
		devices: delivery.NewService(tx),
		me:      principal("a"),
		other:   principal("b"),
		pusher:  &fakePusher{fail: map[string]error{}},
		mailer:  &fakeMailer{},
		prefs: values{
			delivery.LanguagePreference:   `"system"`,
			delivery.QuietHoursPreference: `false`,
			delivery.QuietFromPreference:  `1320`,
			delivery.QuietToPreference:    `420`,
			delivery.TimeZonePreference:   `""`,
		},
		settings: values{
			delivery.LanguageSetting: `"ru"`,
			delivery.TimeZoneSetting: `"Asia/Almaty"`,
		},
	}
}

func (f *fixture) send(t *testing.T, users ...string) {
	t.Helper()
	catalog, err := notify.NewCatalog(deadline)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	err = f.tx.Run(context.Background(), f.me.TenantID, func(tx pgx.Tx) error {
		return notify.NewNotifier(catalog).Send(context.Background(), tx, f.me.TenantID, notify.Message{
			Type: deadline.Key, UserIDs: users, Params: map[string]any{"days": 3}, Link: "/study/1",
		})
	})
	if err != nil {
		t.Fatalf("send: %v", err)
	}
}

func (f *fixture) dispatch(t *testing.T, at time.Time) {
	t.Helper()
	d := delivery.NewDispatcher(f.pool, delivery.DispatcherConfig{
		Texts:       delivery.NewTexts(deadline),
		Preferences: preferences{f.prefs},
		Settings:    settings{f.settings},
		Directory:   identity.Store{},
		Pusher:      f.pusher,
		Mailer:      f.mailer,
		AppURL:      "https://campus.example",
		Every:       time.Hour,
	}).WithClock(func() time.Time { return at })
	if err := d.DispatchAll(context.Background()); err != nil {
		t.Fatalf("dispatch: %v", err)
	}
}

func (f *fixture) pending(t *testing.T) int {
	t.Helper()
	var n int
	err := f.tx.Run(context.Background(), f.me.TenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(context.Background(), `SELECT count(*) FROM core.deliveries`).Scan(&n)
	})
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

func (f *fixture) register(t *testing.T, p auth.Principal, token, language string) delivery.Device {
	t.Helper()
	d, err := f.devices.Register(context.Background(), p, delivery.DeviceInput{Platform: "android", Token: token, Language: language})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	return d
}

func TestPushSpeaksDeviceLanguage(t *testing.T) {
	f := setup(t)
	f.register(t, f.me, "phone", "kk-KZ")
	f.send(t, f.me.UserID)

	f.dispatch(t, time.Now())

	if len(f.pusher.sent) != 1 {
		t.Fatalf("sent = %+v", f.pusher.sent)
	}
	got := f.pusher.sent[0]
	if got.token != "phone" || got.msg.Title != "Тапсыру мерзімі" || got.msg.Body != "3 күн қалды" || got.msg.Data["link"] != "/study/1" {
		t.Fatalf("push = %+v", got)
	}
	if f.pending(t) != 0 {
		t.Fatal("deliveries left in the queue")
	}
}

func TestEmailUsesChosenLanguageAndLink(t *testing.T) {
	f := setup(t)
	f.prefs[delivery.LanguagePreference] = `"en"`
	f.send(t, f.me.UserID)

	f.dispatch(t, time.Now())

	if len(f.mailer.sent) != 1 {
		t.Fatalf("sent = %+v", f.mailer.sent)
	}
	m := f.mailer.sent[0]
	if m.To != "a@uni.example" || m.Subject != "Deadline" || !strings.Contains(m.Body, "https://campus.example/study/1") {
		t.Fatalf("mail = %+v", m)
	}
}

func TestQuietHoursPostponePush(t *testing.T) {
	f := setup(t)
	f.prefs[delivery.QuietHoursPreference] = `true`
	f.register(t, f.me, "phone", "ru")
	f.send(t, f.me.UserID)
	almaty, _ := time.LoadLocation("Asia/Almaty")
	night := time.Date(2026, 11, 12, 23, 30, 0, 0, almaty)

	f.dispatch(t, night)

	if len(f.pusher.sent) != 0 {
		t.Fatalf("pushed during quiet hours: %+v", f.pusher.sent)
	}
	if len(f.mailer.sent) != 1 {
		t.Fatalf("email should not wait for quiet hours: %+v", f.mailer.sent)
	}
	var next time.Time
	_ = f.tx.Run(context.Background(), f.me.TenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(context.Background(), `SELECT next_at FROM core.deliveries WHERE channel = 'push'`).Scan(&next)
	})
	if want := time.Date(2026, 11, 13, 7, 0, 0, 0, almaty); !next.Equal(want) {
		t.Fatalf("next_at = %v, want %v", next, want)
	}
}

func TestUnregisteredDeviceIsForgotten(t *testing.T) {
	f := setup(t)
	f.register(t, f.me, "old", "ru")
	f.register(t, f.me, "new", "ru")
	f.pusher.fail["old"] = fcm.ErrUnregistered
	f.send(t, f.me.UserID)

	f.dispatch(t, time.Now())

	list, _ := f.devices.Devices(context.Background(), f.me)
	if len(list) != 1 || len(f.pusher.sent) != 1 || f.pusher.sent[0].token != "new" {
		t.Fatalf("devices = %+v, sent = %+v", list, f.pusher.sent)
	}
}

func TestFailedEmailRetriesThenGivesUp(t *testing.T) {
	f := setup(t)
	f.mailer.err = errors.New("smtp down")
	f.send(t, f.me.UserID)

	at := time.Now()
	f.dispatch(t, at)
	if f.pending(t) != 1 {
		t.Fatal("failed email was not kept for retry")
	}
	for range 4 {
		at = at.Add(24 * time.Hour)
		_ = f.tx.Run(context.Background(), f.me.TenantID, func(tx pgx.Tx) error {
			_, err := tx.Exec(context.Background(), `UPDATE core.deliveries SET next_at = now()`)
			return err
		})
		f.dispatch(t, at)
	}
	if f.pending(t) != 0 {
		t.Fatal("email kept after five attempts")
	}
}

func TestDeviceMovesToNewOwner(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	mine := f.register(t, f.me, "shared-phone", "ru")
	theirs := f.register(t, f.other, "shared-phone", "ru")

	if mine.ID != theirs.ID {
		t.Fatalf("token created a second device")
	}
	if list, _ := f.devices.Devices(ctx, f.me); len(list) != 0 {
		t.Fatalf("old owner still has device: %+v", list)
	}
	if err := f.devices.Remove(ctx, f.me, theirs.ID); !errors.Is(err, delivery.ErrDeviceNotFound) {
		t.Fatalf("remove someone else's device err = %v", err)
	}
	if err := f.devices.Remove(ctx, f.other, theirs.ID); err != nil {
		t.Fatalf("remove: %v", err)
	}
}

func TestInvalidDeviceIsRejected(t *testing.T) {
	f := setup(t)
	_, err := f.devices.Register(context.Background(), f.me, delivery.DeviceInput{Platform: "symbian", Token: "x"})
	if !errors.Is(err, delivery.ErrInvalidDevice) {
		t.Fatalf("err = %v", err)
	}
}
