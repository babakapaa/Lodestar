package files_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type sharedLink struct {
	file   string
	reader string
}

func (l *sharedLink) Readable(_ context.Context, _ pgx.Tx, p auth.Principal, fileID string) (bool, error) {
	return fileID == l.file && p.UserID == l.reader, nil
}

func (l *sharedLink) InUse(_ context.Context, _ pgx.Tx, fileID string) (bool, error) {
	return fileID == l.file, nil
}

func TestLinkedFileIsReadableAndKept(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	link := &sharedLink{}
	f.service = files.NewService(database.NewTenantTx(f.pool), f.storage, fixedLimit{mb: 1}, link)
	file := f.upload(t, 10)

	if _, err := f.service.Download(ctx, f.other, file.ID); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("stranger download: err = %v", err)
	}
	link.file, link.reader = file.ID, f.other.UserID
	if _, err := f.service.Download(ctx, f.other, file.ID); err != nil {
		t.Fatalf("linked download: %v", err)
	}
	if err := f.service.Delete(ctx, f.other, file.ID); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("reader deletes: err = %v", err)
	}
	if err := f.service.Delete(ctx, f.me, file.ID); !errors.Is(err, files.ErrInUse) {
		t.Fatalf("owner deletes linked: err = %v", err)
	}
	link.file = ""
	if err := f.service.Delete(ctx, f.me, file.ID); err != nil {
		t.Fatalf("owner deletes free file: %v", err)
	}
}
