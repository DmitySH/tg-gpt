package main

import (
	"context"

	"github.com/DmitySH/tg-gpt/internal/app"
	"github.com/DmitySH/tg-gpt/internal/pkg/closer"
	"github.com/DmitySH/tg-gpt/internal/pkg/config"
	"github.com/DmitySH/tg-gpt/internal/pkg/loggy"
)

func main() {
	ctx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	cfg := app.Config{}
	err := config.ReadAndParseYAML(&cfg)
	if err != nil {
		loggy.Fatal("can't read and parse config:", err)
	}

	closer.SetShutdownTimeout(config.Duration("app.graceful_shutdown_timeout"))

	a := app.NewApp(cfg)
	if a.Run(ctx) != nil {
		loggy.Fatalf("failed to run app: %v", err)
	}
}
