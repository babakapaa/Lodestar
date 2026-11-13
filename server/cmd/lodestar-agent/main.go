package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/babakapaa/lodestar/server/internal/modules/sandbox/agent"
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
	gpus, err := strconv.Atoi(env("LODESTAR_AGENT_GPUS", "1"))
	if err != nil {
		return fmt.Errorf("LODESTAR_AGENT_GPUS: %w", err)
	}
	every, err := time.ParseDuration(env("LODESTAR_AGENT_EVERY", "10s"))
	if err != nil {
		return fmt.Errorf("LODESTAR_AGENT_EVERY: %w", err)
	}
	cfg := agent.Config{
		Server:  os.Getenv("LODESTAR_SERVER"),
		Token:   os.Getenv("LODESTAR_NODE_TOKEN"),
		GPUs:    gpus,
		Host:    os.Getenv("LODESTAR_AGENT_HOST"),
		Scheme:  env("LODESTAR_AGENT_SCHEME", "http"),
		DataDir: env("LODESTAR_AGENT_DATA", "/var/lib/lodestar-agent"),
		Every:   every,
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	docker, err := agent.NewDocker(os.Getenv("DOCKER_HOST"))
	if err != nil {
		return err
	}
	agent.New(cfg, docker).Run(ctx)
	return nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
