package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

func New(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      15 * time.Minute,
		IdleTimeout:       120 * time.Second,
	}
}

func Run(ctx context.Context, srv *http.Server, shutdownTimeout time.Duration) error {
	failed := make(chan error, 1)
	go func() {
		slog.Info("http server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			failed <- err
		}
		close(failed)
	}()

	select {
	case err := <-failed:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	slog.Info("http server stopping")
	return srv.Shutdown(shutdownCtx)
}
