package structure

import (
	"encoding/csv"
	"fmt"
	"io"
	"slices"
	"strings"
)

type csvRow struct {
	line   int
	values map[string]string
}

func readTable(r io.Reader, columns ...string) ([]csvRow, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBadImport, err)
	}
	for i := range header {
		header[i] = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(header[i], "\ufeff")))
	}
	for _, c := range columns {
		if !slices.Contains(header, c) {
			return nil, fmt.Errorf("%w: column %s is missing", ErrBadImport, c)
		}
	}
	var rows []csvRow
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return rows, nil
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrBadImport, err)
		}
		line, _ := reader.FieldPos(0)
		row := csvRow{line: line, values: map[string]string{}}
		for i, v := range record {
			if i < len(header) {
				row.values[header[i]] = strings.TrimSpace(v)
			}
		}
		rows = append(rows, row)
	}
}
