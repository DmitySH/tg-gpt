package service

import (
	"context"
	"github.com/DmitySH/tg-gpt/internal/pkg/loggy"

	"github.com/DmitySH/tg-gpt/internal/domain"
	"github.com/sourcegraph/conc/pool"
)

type OnCommandUsecase interface {
	HandleBotCommand(ctx context.Context, update domain.TGUpdate)
}

type OnBasicMessageUsecase interface {
	HandleBasicMessage(ctx context.Context, update domain.TGUpdate)
}

type UpdateProcessor struct {
	cfg    UpdateProcessorConfig
	p      *pool.Pool
	stopCh chan struct{}

	onCommandUsecase      OnCommandUsecase
	onBasicMessageUsecase OnBasicMessageUsecase
}

func NewUpdateProcessor(cfg UpdateProcessorConfig,
	onCommandUsecase OnCommandUsecase) *UpdateProcessor {

	return &UpdateProcessor{
		cfg:              cfg,
		stopCh:           make(chan struct{}),
		p:                pool.New().WithMaxGoroutines(cfg.WorkersCount),
		onCommandUsecase: onCommandUsecase,
	}
}

func (u *UpdateProcessor) ProcessUpdate(update domain.TGUpdate) {
	ctx := context.Background()

	select {
	case <-u.stopCh:
		return
	default:
	}

	u.p.Go(func() {
		defer func() {
			if r := recover(); r != nil {
				loggy.Errorln("recovered from panic", r)
			}
		}()

		u.processUpdate(ctx, update)
	})
}

func (u *UpdateProcessor) processUpdate(ctx context.Context, update domain.TGUpdate) {
	if update.Message != nil && update.Message.Command() != "" {
		u.onCommandUsecase.HandleBotCommand(ctx, update)
		return
	}

	if update.Message != nil {
		u.onBasicMessageUsecase.HandleBasicMessage(ctx, update)
	}
}

func (u *UpdateProcessor) Stop() {
	close(u.stopCh)
	u.p.Wait()
}
