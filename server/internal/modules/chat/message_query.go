package chat

import (
	"net/url"
	"strconv"
)

const (
	pageSize    = 50
	maxPageSize = 200
)

type MessageQuery struct {
	Before *int64
	Since  *int64
	Limit  int
}

func ParseMessageQuery(q url.Values) (MessageQuery, error) {
	var m MessageQuery
	var err error
	if m.Before, err = optionalInt(q.Get("before")); err != nil {
		return MessageQuery{}, err
	}
	if m.Since, err = optionalInt(q.Get("since")); err != nil {
		return MessageQuery{}, err
	}
	if m.Before != nil && m.Since != nil {
		return MessageQuery{}, ErrInvalidMessage
	}
	m.Limit = pageSize
	if raw := q.Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > maxPageSize {
			return MessageQuery{}, ErrInvalidMessage
		}
		m.Limit = n
	}
	return m, nil
}

func optionalInt(raw string) (*int64, error) {
	if raw == "" {
		return nil, nil
	}
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || n < 0 {
		return nil, ErrInvalidMessage
	}
	return &n, nil
}
