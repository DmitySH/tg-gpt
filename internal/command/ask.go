package command

import (
	"github.com/DmitySH/tg-gpt/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AskCommand struct {
	api *tgbotapi.BotAPI
}

func NewAskCommand(api *tgbotapi.BotAPI) AskCommand {
	return AskCommand{
		api: api,
	}
}

func (c AskCommand) String() string {
	return "ask"
}

func (c AskCommand) Process(update domain.TGUpdate) error {
	return nil
}
