package tgbot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	b *tgbotapi.BotAPI
}

func NewBot() Bot {
	bot, err := tgbotapi.NewBotAPI("6863229346:AAFI7GEN40qFDRu4fx7itO7ueLE7gUICVTQ")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	return Bot{b: bot}
}
