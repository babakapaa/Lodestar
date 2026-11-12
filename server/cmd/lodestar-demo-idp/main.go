package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/babakapaa/lodestar/server/internal/app/demo"
	"github.com/babakapaa/lodestar/server/internal/platform/demoidp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	issuer := strings.TrimSuffix(env("DEMO_IDP_ISSUER", "http://localhost:8090"), "/")
	secret := os.Getenv("DEMO_IDP_CLIENT_SECRET")
	if secret == "" {
		return errors.New("DEMO_IDP_CLIENT_SECRET is required")
	}
	var redirects []string
	for item := range strings.SplitSeq(os.Getenv("DEMO_IDP_REDIRECTS"), ",") {
		if item = strings.TrimSpace(item); item != "" {
			redirects = append(redirects, item)
		}
	}
	if len(redirects) == 0 {
		return errors.New("DEMO_IDP_REDIRECTS is required: the Lodestar callback printed by lodestarctl demo seed")
	}
	provider, err := demoidp.New(demoidp.Config{
		Issuer:       issuer,
		ClientID:     env("DEMO_IDP_CLIENT_ID", "lodestar-demo"),
		ClientSecret: secret,
		Redirects:    redirects,
		Users:        demo.People(),
	})
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:              env("DEMO_IDP_ADDR", ":8090"),
		Handler:           provider.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	slog.Info("demo identity provider listening", "addr", server.Addr, "issuer", issuer)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
