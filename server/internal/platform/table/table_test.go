package table_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/platform/table"
)

func TestRoundTrip(t *testing.T) {
	for _, f := range []table.Format{table.CSV, table.XLSX} {
		t.Run(string(f), func(t *testing.T) {
			var buf bytes.Buffer
			err := table.Write(&buf, f, []string{"Email", "Имя", "Баллы"}, [][]any{
				{"a@kaznu.kz", "Айгерим", 8.5},
				{"", "", nil},
				{"b@kaznu.kz", "Данияр, мл.", 10},
			})
			if err != nil {
				t.Fatalf("write: %v", err)
			}
			if table.Detect(buf.Bytes()) != f {
				t.Fatalf("detect = %s", table.Detect(buf.Bytes()))
			}
			got, err := table.Read(&buf)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			if strings.Join(got.Header, "|") != "email|имя|баллы" || len(got.Rows) != 2 {
				t.Fatalf("table = %+v", got)
			}
			if got.Rows[1].Line != 4 || got.Rows[1].Values[1] != "Данияр, мл." || got.Rows[0].Values[2] != "8.5" {
				t.Fatalf("rows = %+v", got.Rows)
			}
		})
	}
}

func TestReadsSemicolonCSV(t *testing.T) {
	got, err := table.Read(strings.NewReader("email;name\nx@kaznu.kz;Икс\n"))
	if err != nil || len(got.Rows) != 1 || got.Rows[0].Values[1] != "Икс" {
		t.Fatalf("table = %+v, %v", got, err)
	}
}

func TestRejectsEmpty(t *testing.T) {
	if _, err := table.Read(strings.NewReader("")); !errors.Is(err, table.ErrMalformed) {
		t.Fatalf("err = %v", err)
	}
	if _, err := table.Read(bytes.NewReader([]byte("PK\x03\x04broken"))); !errors.Is(err, table.ErrMalformed) {
		t.Fatalf("broken xlsx err = %v", err)
	}
}
