package live

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

const (
	pingEvery    = 30 * time.Second
	writeTimeout = 10 * time.Second
)

type StreamHandler struct {
	hub     *Hub
	backlog Backlog
	origins []string
}

func (h StreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	after, ok := afterParam(r)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_request")
		return
	}
	rc := http.NewResponseController(w)
	_ = rc.SetReadDeadline(time.Time{})
	_ = rc.SetWriteDeadline(time.Time{})
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: h.origins})
	if err != nil {
		return
	}
	defer func() { _ = conn.CloseNow() }()
	ctx := conn.CloseRead(r.Context())
	p := principal(r)
	if err := h.serve(ctx, conn, p, after); err != nil && !errors.Is(err, context.Canceled) {
		slog.Debug("live stream ended", "error", err)
	}
}

func (h StreamHandler) serve(ctx context.Context, conn *websocket.Conn, p auth.Principal, after int64) error {
	sub := h.hub.subscribe(p.TenantID, p.UserID)
	defer h.hub.unsubscribe(p.TenantID, p.UserID, sub)

	events, more, err := h.backlog.After(ctx, p, after)
	if err != nil {
		return err
	}
	if more {
		return h.resync(ctx, conn, p)
	}
	last := after
	for _, e := range events {
		if err := write(ctx, conn, e); err != nil {
			return err
		}
		last = e.ID
	}
	ping := time.NewTicker(pingEvery)
	defer ping.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-h.hub.done:
			return conn.Close(websocket.StatusGoingAway, "server stopping")
		case <-sub.lagged:
			return h.resync(ctx, conn, p)
		case e := <-sub.events:
			if e.ID <= last {
				continue
			}
			if err := write(ctx, conn, e); err != nil {
				return err
			}
			last = e.ID
		case <-ping.C:
			pingCtx, cancel := context.WithTimeout(ctx, writeTimeout)
			err := conn.Ping(pingCtx)
			cancel()
			if err != nil {
				return err
			}
		}
	}
}

func (h StreamHandler) resync(ctx context.Context, conn *websocket.Conn, p auth.Principal) error {
	latest, err := h.backlog.Latest(ctx, p)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(Resync{Latest: latest})
	if err != nil {
		return err
	}
	if err := write(ctx, conn, Event{Type: ResyncType, Payload: payload, At: time.Now()}); err != nil {
		return err
	}
	return conn.Close(websocket.StatusNormalClosure, ResyncType)
}

func write(ctx context.Context, conn *websocket.Conn, e Event) error {
	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	return wsjson.Write(ctx, conn, e)
}
