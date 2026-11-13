package structure

import (
	"fmt"
	"io"
	"slices"

	"github.com/babakapaa/lodestar/server/internal/platform/table"
)

type importRow struct {
	line   int
	values map[string]string
}

func readTable(r io.Reader, columns ...string) ([]importRow, error) {
	t, err := table.Read(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBadImport, err)
	}
	for _, c := range columns {
		if !slices.Contains(t.Header, c) {
			return nil, fmt.Errorf("%w: column %s is missing", ErrBadImport, c)
		}
	}
	rows := make([]importRow, 0, len(t.Rows))
	for _, r := range t.Rows {
		row := importRow{line: r.Line, values: map[string]string{}}
		for i, v := range r.Values {
			if i < len(t.Header) {
				row.values[t.Header[i]] = v
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}
