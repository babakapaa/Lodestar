package chat

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const (
	transcribeEvery  = time.Minute
	transcriptRetry  = time.Hour
	transcriptTries  = 3
	transcriptFailed = "failed"
	transcriptDone   = "done"
	transcriptType   = "text/plain; charset=utf-8"
)

type TranscriptWorker struct {
	pool       *pgxpool.Pool
	tenants    database.TenantTx
	k          module.Kernel
	records    RecordingService
	recordings RecordingStore
	files      files.FileStore
	every      time.Duration
}

func NewTranscriptWorker(k module.Kernel) TranscriptWorker {
	return TranscriptWorker{pool: k.Pool, tenants: k.Tenants, k: k, records: NewRecordingService(k), every: transcribeEvery}
}

func (x TranscriptWorker) Run(ctx context.Context) {
	if x.k.Transcriber == nil || x.k.Objects == nil {
		return
	}
	ticker := time.NewTicker(x.every)
	defer ticker.Stop()
	for {
		if err := x.RunAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("transcribe call recordings", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (x TranscriptWorker) RunAll(ctx context.Context) error {
	rows, err := x.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	tenants, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range tenants {
		if err := x.RunTenant(ctx, tenantID); err != nil {
			return err
		}
	}
	return nil
}

func (x TranscriptWorker) RunTenant(ctx context.Context, tenantID string) error {
	enabled, err := x.k.Settings.ModuleEnabled(ctx, tenantID, Name)
	if err != nil || !enabled {
		return err
	}
	for {
		var r Recording
		var tries int
		var found bool
		var f files.File
		err := x.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			var err error
			if r, tries, found, err = x.recordings.ClaimTranscript(ctx, tx, transcriptRetry); err != nil || !found {
				return err
			}
			if r.FileID == nil || r.StartedBy == nil {
				found = false
				return x.recordings.Transcribed(ctx, tx, r.EgressID, transcriptFailed)
			}
			f, err = x.files.ByID(ctx, tx, *r.FileID, false)
			return err
		})
		if err != nil {
			return err
		}
		if !found {
			return nil
		}
		text, err := x.transcribe(ctx, tenantID, f)
		if err != nil {
			slog.Warn("transcribe call recording", "recording", r.EgressID, "try", tries, "error", err)
			if tries < transcriptTries {
				continue
			}
			if err := x.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
				return x.recordings.Transcribed(ctx, tx, r.EgressID, transcriptFailed)
			}); err != nil {
				return err
			}
			continue
		}
		if err := x.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			return x.publish(ctx, tx, tenantID, r, f, text)
		}); err != nil {
			return err
		}
	}
}

func (x TranscriptWorker) transcribe(ctx context.Context, tenantID string, f files.File) (string, error) {
	audio, err := x.k.Objects.Open(ctx, files.ObjectKey(tenantID, f.ID))
	if err != nil {
		return "", err
	}
	defer func() { _ = audio.Close() }()
	text, err := x.k.Transcriber.Transcribe(ctx, audio, f.Name)
	return strings.TrimSpace(text), err
}

func (x TranscriptWorker) publish(ctx context.Context, tx pgx.Tx, tenantID string, r Recording, recording files.File, text string) error {
	if text == "" {
		return x.recordings.Transcribed(ctx, tx, r.EgressID, transcriptDone)
	}
	size := int64(len(text))
	f, err := x.files.Create(ctx, tx, *r.StartedBy, files.FileInput{Name: transcriptFile(recording.Name), ContentType: transcriptType, Size: size})
	if err != nil {
		return err
	}
	if err := x.k.Objects.Put(ctx, files.ObjectKey(tenantID, f.ID), transcriptType, strings.NewReader(text), size); err != nil {
		return err
	}
	if _, err := x.files.Finish(ctx, tx, f.ID, size); err != nil {
		return err
	}
	if err := x.records.post(ctx, tx, tenantID, r, f.ID, "call-transcript-"+r.EgressID); err != nil {
		return err
	}
	return x.recordings.Transcribed(ctx, tx, r.EgressID, transcriptDone)
}
