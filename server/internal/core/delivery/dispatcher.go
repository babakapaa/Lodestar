package delivery

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/fcm"
	"github.com/babakapaa/lodestar/server/internal/platform/mail"
)

const (
	batchSize   = 100
	maxAttempts = 5
)

var backoff = []time.Duration{time.Minute, 5 * time.Minute, 30 * time.Minute, 2 * time.Hour}

type Dispatcher struct {
	pool       *pgxpool.Pool
	tenants    database.TenantTx
	texts      Texts
	recipients recipientReader
	directory  Directory
	pusher     Pusher
	mailer     Mailer
	appURL     string
	every      time.Duration
	now        func() time.Time
	jobs       JobStore
	devices    DeviceStore
}

type DispatcherConfig struct {
	Texts       Texts
	Preferences Preferences
	Settings    Settings
	Directory   Directory
	Pusher      Pusher
	Mailer      Mailer
	AppURL      string
	Every       time.Duration
}

func NewDispatcher(pool *pgxpool.Pool, cfg DispatcherConfig) Dispatcher {
	return Dispatcher{
		pool:       pool,
		tenants:    database.NewTenantTx(pool),
		texts:      cfg.Texts,
		recipients: recipientReader{preferences: cfg.Preferences, settings: cfg.Settings},
		directory:  cfg.Directory,
		pusher:     cfg.Pusher,
		mailer:     cfg.Mailer,
		appURL:     cfg.AppURL,
		every:      cfg.Every,
		now:        time.Now,
	}
}

func (d Dispatcher) WithClock(now func() time.Time) Dispatcher {
	d.now = now
	return d
}

func (d Dispatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(d.every)
	defer ticker.Stop()
	for {
		if err := d.DispatchAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("dispatch deliveries", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (d Dispatcher) DispatchAll(ctx context.Context) error {
	rows, err := d.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := d.tenants.Run(ctx, id, func(tx pgx.Tx) error { return d.dispatch(ctx, tx) }); err != nil {
			return err
		}
	}
	return nil
}

func (d Dispatcher) dispatch(ctx context.Context, tx pgx.Tx) error {
	jobs, err := d.jobs.Claim(ctx, tx, batchSize)
	if err != nil {
		return err
	}
	for _, j := range jobs {
		if err := d.handle(ctx, tx, j); err != nil {
			return err
		}
	}
	return nil
}

func (d Dispatcher) handle(ctx context.Context, tx pgx.Tx, j job) error {
	r, err := d.recipients.read(ctx, tx, j.userID)
	if err != nil {
		return err
	}
	switch j.channel {
	case ChannelPush:
		if until, quiet := r.quiet.until(d.now()); quiet {
			return d.jobs.Postpone(ctx, tx, j.id, until)
		}
		err = d.push(ctx, tx, j, r)
	case ChannelEmail:
		err = d.email(ctx, tx, j, r)
	}
	var failure sendFailure
	if err != nil && !errors.As(err, &failure) {
		return err
	}
	if err == nil {
		return d.jobs.Done(ctx, tx, j.id)
	}
	if j.attempts+1 >= maxAttempts {
		slog.Warn("delivery dropped after retries", "channel", j.channel, "type", j.item.Type, "error", failure)
		return d.jobs.Done(ctx, tx, j.id)
	}
	return d.jobs.Retry(ctx, tx, j.id, d.now().Add(backoff[min(j.attempts, len(backoff)-1)]), failure.Error())
}

func (d Dispatcher) push(ctx context.Context, tx pgx.Tx, j job, r recipient) error {
	if d.pusher == nil {
		return nil
	}
	targets, err := d.devices.Targets(ctx, tx, j.userID)
	if err != nil {
		return err
	}
	var failed error
	delivered := len(targets) == 0
	for _, t := range targets {
		text, ok := d.texts.Render(j.item, r.language, t.language, r.tenantLanguage)
		if !ok {
			return nil
		}
		msg := fcm.Message{Title: text.Title, Body: text.Body, Data: map[string]string{"type": j.item.Type, "link": j.item.Link}}
		switch err := d.pusher.Push(ctx, t.token, msg); {
		case err == nil:
			delivered = true
		case errors.Is(err, fcm.ErrUnregistered):
			if err := d.devices.Forget(ctx, tx, t.id); err != nil {
				return err
			}
		default:
			failed = err
		}
	}
	if delivered || failed == nil {
		return nil
	}
	return sendFailure{err: failed}
}

func (d Dispatcher) email(ctx context.Context, tx pgx.Tx, j job, r recipient) error {
	if d.mailer == nil {
		return nil
	}
	address, err := d.directory.Email(ctx, tx, j.userID)
	if err != nil || address == "" {
		return err
	}
	targets, err := d.devices.Targets(ctx, tx, j.userID)
	if err != nil {
		return err
	}
	deviceLanguage := ""
	if len(targets) > 0 {
		deviceLanguage = targets[0].language
	}
	text, ok := d.texts.Render(j.item, r.language, r.tenantLanguage, deviceLanguage)
	if !ok {
		return nil
	}
	body := text.Body
	if j.item.Link != "" {
		body += "\n\n" + d.appURL + j.item.Link
	}
	if err := d.mailer.Send(ctx, mail.Message{To: address, Subject: text.Title, Body: body}); err != nil {
		return sendFailure{err: err}
	}
	return nil
}
