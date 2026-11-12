package study_test

import (
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestAttachedFilesAreInUse(t *testing.T) {
	w := newWorld(t)
	c := w.course(nil, "Химия")
	attached := w.file(w.admin, "a.txt")
	free := w.file(w.admin, "b.txt")
	if _, err := w.work.Create(w.ctx, w.admin, c.ID, study.AssignmentInput{Title: "Опыт", Files: []string{attached}}); err != nil {
		t.Fatalf("create: %v", err)
	}
	link := study.Module{}.FileLinks(module.Kernel{Tenants: w.tx, Access: w.access})[0]
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		used, err := link.InUse(w.ctx, tx, attached)
		if err != nil || !used {
			t.Errorf("attached in use = %v, %v", used, err)
		}
		used, err = link.InUse(w.ctx, tx, free)
		if err != nil || used {
			t.Errorf("free in use = %v, %v", used, err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
