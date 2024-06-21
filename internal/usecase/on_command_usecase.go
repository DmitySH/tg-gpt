package usecase

import (
	"fmt"

	"github.com/DmitySH/tg-gpt/internal/command"
	"github.com/DmitySH/tg-gpt/internal/domain"
	"github.com/DmitySH/tg-gpt/internal/pkg/loggy"
)

type OnCommandUsecase struct {
	askCommand command.AskCommand
}

func NewOnCommandUsecase(askCommand command.AskCommand) OnCommandUsecase {
	return OnCommandUsecase{
		askCommand: askCommand,
	}
}

func (u OnCommandUsecase) HandleBotCommand(update domain.TGUpdate) {
	var err error
	switch update.Message.Command() {
	case u.askCommand.String():
		err = u.askCommand.Process(update)
		if err != nil {
			err = fmt.Errorf("can't process ask command: %w", err)
		}
	}

	if err != nil {
		loggy.Errorln(err)
	}
}
