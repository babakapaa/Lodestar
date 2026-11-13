package table

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

const sheetName = "Lodestar"

func Write(w io.Writer, f Format, header []string, rows [][]any) error {
	if f == XLSX {
		return writeXLSX(w, header, rows)
	}
	return writeCSV(w, header, rows)
}

func writeCSV(w io.Writer, header []string, rows [][]any) error {
	if _, err := io.WriteString(w, "\ufeff"); err != nil {
		return err
	}
	out := csv.NewWriter(w)
	if err := out.Write(header); err != nil {
		return err
	}
	for _, row := range rows {
		record := make([]string, len(row))
		for i, v := range row {
			if v != nil {
				record[i] = fmt.Sprint(v)
			}
		}
		if err := out.Write(record); err != nil {
			return err
		}
	}
	out.Flush()
	return out.Error()
}

func writeXLSX(w io.Writer, header []string, rows [][]any) error {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	if err := f.SetSheetName(f.GetSheetName(0), sheetName); err != nil {
		return err
	}
	stream, err := f.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}
	head := make([]any, len(header))
	for i, h := range header {
		head[i] = h
	}
	if err := stream.SetRow("A1", head); err != nil {
		return err
	}
	for i, row := range rows {
		cell, err := excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			return err
		}
		if err := stream.SetRow(cell, row); err != nil {
			return err
		}
	}
	if err := stream.Flush(); err != nil {
		return err
	}
	_, err = f.WriteTo(w)
	return err
}
