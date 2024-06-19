package app

import (
	"context"
	"time"

	"github.com/DmitySH/tg-gpt/internal/pkg/closer"
	"github.com/DmitySH/tg-gpt/internal/pkg/loggy"
)

type App struct {
	cfg Config
}

func NewApp(config Config) *App {
	return &App{
		cfg: config,
	}
}

func (a *App) Run(ctx context.Context) error {
	closer.Add(closer.NoErrAdapter(func() {
		time.Sleep(time.Second * 3)
	}))

	closer.Wait()
	loggy.Sync()

	return nil
}
