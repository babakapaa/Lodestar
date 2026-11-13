package table

import "bytes"

type Format string

const (
	CSV  Format = "csv"
	XLSX Format = "xlsx"
)

var zipMagic = []byte("PK\x03\x04")

func ParseFormat(value string) (Format, bool) {
	switch Format(value) {
	case "", CSV:
		return CSV, true
	case XLSX:
		return XLSX, true
	default:
		return "", false
	}
}

func Detect(data []byte) Format {
	if bytes.HasPrefix(data, zipMagic) {
		return XLSX
	}
	return CSV
}

func (f Format) ContentType() string {
	if f == XLSX {
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
	return "text/csv; charset=utf-8"
}

func (f Format) Extension() string {
	return "." + string(f)
}
