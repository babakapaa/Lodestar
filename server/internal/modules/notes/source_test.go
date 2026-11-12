package notes_test

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/notes"
)

const sourceID = "5b1c7a9e-2f3d-4c6b-9a8e-1d2f3a4b5c6d"

func TestSourcesAreOwnedAndEncrypted(t *testing.T) {
	w := newWorld(t)
	aigerim := w.user("aigerim")
	erlan := w.user("erlan")
	in := notes.SourceInput{ID: sourceID, SourceFields: notes.SourceFields{
		Kind: "book", Title: "  Основы цитологии ", Authors: []string{"Иванов, А. Б.", " "}, DOI: "https://doi.org/10.1/abc",
	}}

	src, created, err := w.notes.SaveSource(w.ctx, aigerim, in)
	if err != nil || !created {
		t.Fatalf("save: %v %v", created, err)
	}
	if src.Title != "Основы цитологии" || len(src.Authors) != 1 || src.DOI != "10.1/abc" {
		t.Fatalf("normalized: %+v", src)
	}
	if raw := w.raw(`SELECT data FROM notes.sources WHERE id = $1`, sourceID); bytes.Contains(raw, []byte("цитологии")) {
		t.Fatal("source is stored in plain text")
	}
	in.Title = "Основы цитологии, 2-е изд."
	if _, created, err := w.notes.SaveSource(w.ctx, aigerim, in); err != nil || created {
		t.Fatalf("update: %v %v", created, err)
	}
	if _, _, err := w.notes.SaveSource(w.ctx, erlan, in); !errors.Is(err, notes.ErrSourceNotFound) {
		t.Fatalf("someone else overwrote: %v", err)
	}
	list, err := w.notes.Sources(w.ctx, erlan)
	if err != nil || len(list) != 0 {
		t.Fatalf("erlan sees %d: %v", len(list), err)
	}
	list, err = w.notes.Sources(w.ctx, aigerim)
	if err != nil || len(list) != 1 || list[0].Title != "Основы цитологии, 2-е изд." {
		t.Fatalf("list: %+v %v", list, err)
	}
	if err := w.notes.DeleteSource(w.ctx, erlan, sourceID); !errors.Is(err, notes.ErrSourceNotFound) {
		t.Fatalf("someone else deleted: %v", err)
	}
	if err := w.notes.DeleteSource(w.ctx, aigerim, sourceID); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestSourceValidation(t *testing.T) {
	w := newWorld(t)
	p := w.user("aigerim")
	cases := []notes.SourceInput{
		{ID: sourceID, SourceFields: notes.SourceFields{Title: ""}},
		{ID: "x", SourceFields: notes.SourceFields{Title: "a"}},
		{ID: sourceID, SourceFields: notes.SourceFields{Title: "a", Kind: "tweet"}},
		{ID: sourceID, SourceFields: notes.SourceFields{Title: strings.Repeat("я", 1001)}},
	}
	for i, in := range cases {
		if _, _, err := w.notes.SaveSource(w.ctx, p, in); !errors.Is(err, notes.ErrInvalidSource) {
			t.Errorf("case %d: %v", i, err)
		}
	}
	fileID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
	in := notes.SourceInput{ID: sourceID, SourceFields: notes.SourceFields{Title: "a"}, FileID: &fileID}
	if _, _, err := w.notes.SaveSource(w.ctx, p, in); !errors.Is(err, notes.ErrInvalidSource) {
		t.Errorf("unknown file: %v", err)
	}
}

func TestLookupByDOIAndPage(t *testing.T) {
	w := newWorld(t)
	byDOI, err := w.notes.LookupSource(w.ctx, notes.LookupInput{DOI: "doi:10.1000/xyz"})
	if err != nil {
		t.Fatalf("doi: %v", err)
	}
	if byDOI.Kind != "article" || byDOI.Title != "Строение клетки" || byDOI.Year != "2021" ||
		byDOI.Pages != "45–67" || len(byDOI.Authors) != 2 || byDOI.Authors[0] != "Иванов, А. Б." {
		t.Fatalf("doi fields: %+v", byDOI)
	}
	page, err := w.notes.LookupSource(w.ctx, notes.LookupInput{URL: "https://journal.example.org/a/1"})
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	if page.Title != "Митоз и мейоз" || page.Kind != "article" || page.Year != "2019" ||
		page.DOI != "10.1/m" || len(page.Authors) != 2 || page.AccessedAt == "" {
		t.Fatalf("page fields: %+v", page)
	}
	if _, err := w.notes.LookupSource(w.ctx, notes.LookupInput{URL: "file:///etc/passwd"}); !errors.Is(err, notes.ErrInvalidSource) {
		t.Fatalf("file url: %v", err)
	}
	if _, err := w.notes.LookupSource(w.ctx, notes.LookupInput{DOI: "10.404/none"}); !errors.Is(err, notes.ErrLookupFailed) {
		t.Fatalf("missing doi: %v", err)
	}
}

func TestSourcesOverHTTP(t *testing.T) {
	w := newWorld(t)
	p := w.user("aigerim")
	rec := w.serve(p, http.MethodPost, "/v1/notes/sources", `{"id":"`+sourceID+`","kind":"web","title":"Википедия"}`)
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"title":"Википедия"`) {
		t.Fatalf("save: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(p, http.MethodGet, "/v1/notes/sources", ""); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), sourceID) {
		t.Fatalf("list: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(p, http.MethodPost, "/v1/notes/lookup", `{"doi":"10.1000/xyz"}`); rec.Code != http.StatusOK {
		t.Fatalf("lookup: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(p, http.MethodPost, "/v1/notes/lookup", `{}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("empty lookup: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(p, http.MethodDelete, "/v1/notes/sources/"+sourceID, ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(p, http.MethodGet, "/v1/notes", ""); rec.Code != http.StatusOK {
		t.Fatalf("notes list still works: %d", rec.Code)
	}
}
