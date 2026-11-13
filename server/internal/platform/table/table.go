package table

import "errors"

var ErrMalformed = errors.New("table file is malformed")

type Table struct {
	Header []string
	Rows   []Row
}

type Row struct {
	Line   int
	Values []string
}
