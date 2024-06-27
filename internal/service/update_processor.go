package service

import (
	"context"

	"github.com/DmitySH/tg-gpt/internal/domain"
	"github.com/DmitySH/tg-gpt/internal/pkg/loggy"
	"github.com/sourcegraph/conc/pool"
)

type OnCommandUsecase interface {
	HandleBotCommand(ctx context.Context, update domain.TGUpdate) error
}

type OnBasicMessageUsecase interface {
	HandleBasicMessage(ctx context.Context, update domain.TGUpdate) error
}

type UpdateProcessor struct {
	cfg    UpdateProcessorConfig
	p      *pool.Pool
	stopCh chan struct{}

	onCommandUsecase      OnCommandUsecase
	onBasicMessageUsecase OnBasicMessageUsecase
}

func NewUpdateProcessor(cfg UpdateProcessorConfig,
	onCommandUsecase OnCommandUsecase, onBasicMessageUsecase OnBasicMessageUsecase) *UpdateProcessor {

	return &UpdateProcessor{
		cfg:                   cfg,
		stopCh:                make(chan struct{}),
		p:                     pool.New().WithMaxGoroutines(cfg.WorkersCount),
		onCommandUsecase:      onCommandUsecase,
		onBasicMessageUsecase: onBasicMessageUsecase,
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

		err := u.processUpdate(ctx, update)
		if err != nil {
			loggy.Errorln("can't process update:", err)
		}
	})
}

func (u *UpdateProcessor) processUpdate(ctx context.Context, update domain.TGUpdate) error {
	if update.Message != nil && update.Message.Command() != "" {
		return u.onCommandUsecase.HandleBotCommand(ctx, update)
	}

	if update.Message != nil {
		return u.onBasicMessageUsecase.HandleBasicMessage(ctx, update)
	}

	return nil
}

func (u *UpdateProcessor) Stop() {
	close(u.stopCh)
	u.p.Wait()
}
