package league

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	maxDivisions = 10
	maxBoundary  = 1000000
)

type Divisions []Division

func defaultDivisions() Divisions {
	return Divisions{{From: 0}, {From: 150}, {From: 300}}
}

func (d Divisions) normalized() (Divisions, error) {
	if len(d) == 0 || len(d) > maxDivisions || d[0].From != 0 {
		return nil, ErrInvalidRules
	}
	out := make(Divisions, len(d))
	for i, div := range d {
		div.Name = strings.TrimSpace(div.Name)
		if utf8.RuneCountInString(div.Name) > 40 || div.From > maxBoundary || (i > 0 && div.From <= d[i-1].From) {
			return nil, ErrInvalidRules
		}
		out[i] = div
	}
	return out, nil
}

func (d Divisions) at(points int) int {
	level := 0
	for i, div := range d {
		if points >= div.From {
			level = i
		}
	}
	return level
}

func (d Divisions) clamp(level int) int {
	return max(0, min(level, len(d)-1))
}

func (d Divisions) toNext(level, points int) *int {
	if level+1 >= len(d) {
		return nil
	}
	left := max(d[level+1].From-points, 0)
	return &left
}

func (d Divisions) label(level int) string {
	if name := d[d.clamp(level)].Name; name != "" {
		return name
	}
	return strconv.Itoa(level + 1)
}
