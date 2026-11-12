package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const usage = `Использование:
  lodestarctl secret-key
  lodestarctl tenant create -slug kaznu -name "Казахский национальный университет"
  lodestarctl oidc set -tenant kaznu -key microsoft -name "Microsoft" -issuer URL -client-id ID
      секрет клиента берётся из LODESTAR_OIDC_CLIENT_SECRET`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if len(args) == 1 && args[0] == "secret-key" {
		return printSecretKey()
	}
	if len(args) < 2 {
		return fmt.Errorf("%s", usage)
	}
	switch args[0] + " " + args[1] {
	case "tenant create":
		return createTenant(ctx, args[2:])
	case "oidc set":
		return setOIDC(ctx, args[2:])
	default:
		return fmt.Errorf("%s", usage)
	}
}
