package notes_test

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/modules/notes"
)

const noteID = "6f1d2c3b-4a5e-4f60-8a7b-9c0d1e2f3a4b"

func TestTwoPeopleEditOneNote(t *testing.T) {
	w := newWorld(t)
	owner, editor, reader, stranger := w.user("aigerim"), w.user("erlan"), w.user("dana"), w.user("zhan")
	w.create(owner, "Биология/Клетка")

	if err := w.notes.SetMember(w.ctx, owner, noteID, editor.UserID, notes.MemberInput{CanEdit: true}); err != nil {
		t.Fatalf("share: %v", err)
	}
	if err := w.notes.SetMember(w.ctx, owner, noteID, reader.UserID, notes.MemberInput{}); err != nil {
		t.Fatalf("share read: %v", err)
	}
	if err := w.notes.SetMember(w.ctx, editor, noteID, stranger.UserID, notes.MemberInput{}); !errors.Is(err, notes.ErrNotOwner) {
		t.Fatalf("editor shares = %v", err)
	}

	if seq := w.push(owner, "a-1", "offline edit A"); seq != 1 {
		t.Fatalf("first seq = %d", seq)
	}
	if seq := w.push(editor, "b-1", "offline edit B"); seq != 2 {
		t.Fatalf("second seq = %d", seq)
	}
	if seq := w.push(editor, "b-1", "offline edit B"); seq != 2 {
		t.Fatalf("retried push seq = %d", seq)
	}
	if _, _, err := w.notes.Push(w.ctx, reader, noteID, notes.UpdateInput{Key: "c-1", Data: []byte("x")}); !errors.Is(err, notes.ErrReadOnly) {
		t.Fatalf("reader push = %v", err)
	}
	if _, err := w.notes.Changes(w.ctx, stranger, noteID, 0); !errors.Is(err, notes.ErrNoteNotFound) {
		t.Fatalf("stranger read = %v", err)
	}

	c, err := w.notes.Changes(w.ctx, reader, noteID, 1)
	if err != nil || c.Seq != 2 || len(c.Updates) != 1 || string(c.Updates[0].Data) != "offline edit B" {
		t.Fatalf("changes = %+v, %v", c, err)
	}
	list, err := w.notes.List(w.ctx, editor)
	if err != nil || len(list) != 1 || list[0].Role != notes.RoleEditor || list[0].Meta.Path != "Биология/Клетка" {
		t.Fatalf("editor list = %+v, %v", list, err)
	}
	if _, _, err := w.notes.Create(w.ctx, stranger, notes.CreateInput{ID: noteID, Meta: notes.Meta{Path: "x"}}); !errors.Is(err, notes.ErrNoteExists) {
		t.Fatalf("hijack id = %v", err)
	}
}

func TestCompactionKeepsVersionsAndPrunesUpdates(t *testing.T) {
	w := newWorld(t)
	owner := w.user("aigerim")
	w.create(owner, "Клетка")
	for _, key := range []string{"1", "2", "3"} {
		w.push(owner, key, "update "+key)
	}
	if _, err := w.notes.Compact(w.ctx, owner, noteID, notes.SnapshotInput{Seq: 4, Data: []byte("s"), Body: "x"}); !errors.Is(err, notes.ErrStaleSnapshot) {
		t.Fatalf("future snapshot = %v", err)
	}
	v, err := w.notes.Compact(w.ctx, owner, noteID, notes.SnapshotInput{Seq: 2, Data: []byte("state at 2"), Body: "# Клетка"})
	if err != nil {
		t.Fatalf("compact: %v", err)
	}
	if _, err := w.notes.Compact(w.ctx, owner, noteID, notes.SnapshotInput{Seq: 2, Data: []byte("again"), Body: ""}); !errors.Is(err, notes.ErrStaleSnapshot) {
		t.Fatalf("repeated snapshot = %v", err)
	}
	c, err := w.notes.Changes(w.ctx, owner, noteID, 0)
	if err != nil || c.Snapshot == nil || string(c.Snapshot.Data) != "state at 2" || len(c.Updates) != 1 || c.Updates[0].Seq != 3 {
		t.Fatalf("changes after compaction = %+v, %v", c, err)
	}
	c, err = w.notes.Changes(w.ctx, owner, noteID, 2)
	if err != nil || c.Snapshot != nil || len(c.Updates) != 1 {
		t.Fatalf("changes from 2 = %+v, %v", c, err)
	}
	content, err := w.notes.Version(w.ctx, owner, noteID, v.ID)
	if err != nil || content.Body != "# Клетка" || content.Seq != 2 {
		t.Fatalf("version = %+v, %v", content, err)
	}
	if list, err := w.notes.Versions(w.ctx, owner, noteID); err != nil || len(list) != 1 {
		t.Fatalf("versions = %+v, %v", list, err)
	}
}

func TestNotesAreEncryptedAtRest(t *testing.T) {
	w := newWorld(t)
	owner := w.user("aigerim")
	w.create(owner, "Секретная/Тема")
	w.push(owner, "k", "текст заметки")
	if _, err := w.notes.Compact(w.ctx, owner, noteID, notes.SnapshotInput{Seq: 1, Data: []byte("state"), Body: "тело заметки"}); err != nil {
		t.Fatalf("compact: %v", err)
	}
	w.push(owner, "k2", "вторая правка")
	for _, check := range []struct {
		query, plain string
	}{
		{`SELECT meta FROM notes.notes`, "Секретная"},
		{`SELECT data FROM notes.updates`, "вторая"},
		{`SELECT body FROM notes.versions`, "тело"},
	} {
		if raw := w.raw(check.query); bytes.Contains(raw, []byte(check.plain)) {
			t.Errorf("%s stores plain text", check.query)
		}
	}
}

func TestQuotaCountsOwnerBytes(t *testing.T) {
	w := newWorld(t)
	owner, editor := w.user("aigerim"), w.user("erlan")
	w.create(owner, "Большая")
	if err := w.notes.SetMember(w.ctx, owner, noteID, editor.UserID, notes.MemberInput{CanEdit: true}); err != nil {
		t.Fatalf("share: %v", err)
	}
	w.settings.values[notes.QuotaSetting] = 1
	big := strings.Repeat("x", 600<<10)
	w.push(owner, "1", big)
	if _, _, err := w.notes.Push(w.ctx, editor, noteID, notes.UpdateInput{Key: "2", Data: []byte(big)}); !errors.Is(err, notes.ErrQuotaExceeded) {
		t.Fatalf("over quota = %v", err)
	}
}

func TestTrashAndPurge(t *testing.T) {
	w := newWorld(t)
	owner, editor := w.user("aigerim"), w.user("erlan")
	w.create(owner, "Удалить")
	if err := w.notes.SetMember(w.ctx, owner, noteID, editor.UserID, notes.MemberInput{CanEdit: true}); err != nil {
		t.Fatalf("share: %v", err)
	}
	if err := w.notes.SetDeleted(w.ctx, editor, noteID, true); !errors.Is(err, notes.ErrNotOwner) {
		t.Fatalf("editor delete = %v", err)
	}
	if err := w.notes.SetDeleted(w.ctx, owner, noteID, true); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, _, err := w.notes.Push(w.ctx, editor, noteID, notes.UpdateInput{Key: "x", Data: []byte("x")}); !errors.Is(err, notes.ErrNoteDeleted) {
		t.Fatalf("push to deleted = %v", err)
	}
	purger := notes.NewTrashPurger(w.kernel)
	if err := purger.PurgeAll(w.ctx); err != nil {
		t.Fatalf("purge now: %v", err)
	}
	if list, _ := w.notes.List(w.ctx, owner); len(list) != 1 || list[0].DeletedAt == nil {
		t.Fatalf("trash = %+v", list)
	}
	w.settings.values[notes.TrashDaysSetting] = 1
	w.raw(`UPDATE notes.notes SET deleted_at = now() - interval '2 days' RETURNING meta`)
	if err := purger.PurgeAll(w.ctx); err != nil {
		t.Fatalf("purge: %v", err)
	}
	if list, _ := w.notes.List(w.ctx, owner); len(list) != 0 {
		t.Fatalf("after purge = %+v", list)
	}
}

func TestPrivacyExportOpensNotes(t *testing.T) {
	w := newWorld(t)
	owner := w.user("aigerim")
	w.create(owner, "Клетка")
	w.push(owner, "1", "x")
	if _, err := w.notes.Compact(w.ctx, owner, noteID, notes.SnapshotInput{Seq: 1, Data: []byte("s"), Body: "# Ядро"}); err != nil {
		t.Fatalf("compact: %v", err)
	}
	section := notes.Module{}.PrivacySections(w.kernel)[0]
	var out any
	err := w.kernel.Tenants.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		out, err = section.Export(w.ctx, tx, owner.UserID)
		return err
	})
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	list := out.([]map[string]any)
	if len(list) != 1 || list[0]["path"] != "Клетка" || list[0]["body"] != "# Ядро" {
		t.Fatalf("export = %+v", list)
	}
}

func TestNotesOverHTTP(t *testing.T) {
	w := newWorld(t)
	owner := w.user("aigerim")
	rec := w.serve(owner, http.MethodPost, "/v1/notes", `{"id":"`+noteID+`","meta":{"path":"/Биология/Клетка/","tags":["клетка"]}}`)
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"path":"Биология/Клетка"`) {
		t.Fatalf("create: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(owner, http.MethodPost, "/v1/notes", `{"id":"`+noteID+`","meta":{"path":"x"}}`); rec.Code != http.StatusOK {
		t.Fatalf("repeat create: %d %s", rec.Code, rec.Body)
	}
	rec = w.serve(owner, http.MethodPost, "/v1/notes/"+noteID+"/updates", `{"key":"k1","data":"AQID"}`)
	if rec.Code != http.StatusCreated || rec.Body.String() != "{\"seq\":1}\n" {
		t.Fatalf("push: %d %s", rec.Code, rec.Body)
	}
	rec = w.serve(owner, http.MethodGet, "/v1/notes/"+noteID+"/updates?after=0", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"data":"AQID"`) {
		t.Fatalf("changes: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(owner, http.MethodGet, "/v1/notes/"+noteID+"/updates?after=-1", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad after: %d", rec.Code)
	}
	if rec := w.serve(owner, http.MethodPut, "/v1/notes/"+noteID+"/meta", `{"path":""}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("empty path: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(owner, http.MethodDelete, "/v1/notes/"+noteID, ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(owner, http.MethodPost, "/v1/notes/"+noteID+"/restore", ""); rec.Code != http.StatusNoContent {
		t.Fatalf("restore: %d %s", rec.Code, rec.Body)
	}
}
