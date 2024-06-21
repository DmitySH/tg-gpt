package service

import (
	"github.com/DmitySH/tg-gpt/internal/domain"
	"github.com/sourcegraph/conc/pool"
)

type OnCommandUsecase interface {
	HandleBotCommand(update domain.TGUpdate)
}

type UpdateProcessor struct {
	cfg    UpdateProcessorConfig
	p      *pool.Pool
	stopCh chan struct{}

	onCommandUsecase OnCommandUsecase
}

func NewUpdateProcessor(cfg UpdateProcessorConfig,
	onCommandUsecase OnCommandUsecase) UpdateProcessor {

	return UpdateProcessor{
		cfg:              cfg,
		stopCh:           make(chan struct{}),
		p:                pool.New().WithMaxGoroutines(cfg.WorkersCount),
		onCommandUsecase: onCommandUsecase,
	}
}

func (u UpdateProcessor) ProcessUpdate(update domain.TGUpdate) {
	select {
	case <-u.stopCh:
		return
	default:
	}

	u.p.Go(func() {
		if update.Message == nil {
			return
		}

		u.routeUpdate(update)
	})
}

func (u UpdateProcessor) routeUpdate(update domain.TGUpdate) {
	if update.Message.Command() != "" {
		u.onCommandUsecase.HandleBotCommand(update)
	}
}

func (u UpdateProcessor) Stop() {
	close(u.stopCh)
	u.p.Wait()
}
