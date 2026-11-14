package mail

import (
	"context"
	"net"
)

func Ping(ctx context.Context, cfg Config) error {
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", cfg.Addr)
	if err != nil {
		return err
	}
	return conn.Close()
}
