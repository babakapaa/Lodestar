package table

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

func Read(r io.Reader) (Table, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Table{}, err
	}
	var records []Row
	if Detect(data) == XLSX {
		records, err = readXLSX(data)
	} else {
		records, err = readCSV(data)
	}
	if err != nil {
		return Table{}, fmt.Errorf("%w: %w", ErrMalformed, err)
	}
	if len(records) == 0 {
		return Table{}, fmt.Errorf("%w: header is missing", ErrMalformed)
	}
	t := Table{Header: make([]string, len(records[0].Values))}
	for i, name := range records[0].Values {
		t.Header[i] = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(name, "\ufeff")))
	}
	for _, record := range records[1:] {
		empty := true
		for j, v := range record.Values {
			record.Values[j] = strings.TrimSpace(v)
			empty = empty && record.Values[j] == ""
		}
		if !empty {
			t.Rows = append(t.Rows, record)
		}
	}
	return t, nil
}

func readCSV(data []byte) ([]Row, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	if sample := firstLine(data); strings.Count(sample, ";") > strings.Count(sample, ",") {
		reader.Comma = ';'
	}
	var rows []Row
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return rows, nil
		}
		if err != nil {
			return nil, err
		}
		line, _ := reader.FieldPos(0)
		rows = append(rows, Row{Line: line, Values: record})
	}
}

func readXLSX(data []byte) ([]Row, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("workbook has no sheets")
	}
	records, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}
	rows := make([]Row, len(records))
	for i, record := range records {
		rows[i] = Row{Line: i + 1, Values: record}
	}
	return rows, nil
}

func firstLine(data []byte) string {
	line, _, _ := strings.Cut(string(data), "\n")
	return line
}
