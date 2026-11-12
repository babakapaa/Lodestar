package live

import (
	"fmt"
	"strconv"
	"strings"
)

const notifyChannel = "lodestar_live"

type signal struct {
	tenantID string
	userID   string
	eventID  int64
}

func (s signal) encode() string {
	return s.tenantID + ":" + s.userID + ":" + strconv.FormatInt(s.eventID, 10)
}

func decodeSignal(payload string) (signal, error) {
	parts := strings.Split(payload, ":")
	if len(parts) != 3 {
		return signal{}, fmt.Errorf("bad live signal %q", payload)
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return signal{}, fmt.Errorf("bad live signal %q: %w", payload, err)
	}
	return signal{tenantID: parts[0], userID: parts[1], eventID: id}, nil
}
