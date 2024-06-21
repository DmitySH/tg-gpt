package tgbot

import (
	"time"

	"github.com/DmitySH/tg-gpt/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotReader struct {
	cfg Config
	api *tgbotapi.BotAPI

	updateProcessor service.UpdateProcessor
}

func NewBotReader(cfg Config, botAPI *tgbotapi.BotAPI, updateProcessor service.UpdateProcessor) BotReader {
	return BotReader{
		api:             botAPI,
		cfg:             cfg,
		updateProcessor: updateProcessor,
	}
}

func (b BotReader) StartReceivingUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.api.GetUpdatesChan(u)
	firstUpdate, ok := skipOldUpdatesAndGetFirstActual(updates)
	if ok {
		b.updateProcessor.ProcessUpdate(firstUpdate)
	}

	go func() {
		for update := range updates {
			b.updateProcessor.ProcessUpdate(update)
		}
	}()
}

func skipOldUpdatesAndGetFirstActual(updates tgbotapi.UpdatesChannel) (tgbotapi.Update, bool) {
	startTime := time.Now().UTC()
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if time.Unix(int64(update.Message.Date), 0).UTC().After(startTime.Add(-time.Second * 5)) {
			return update, true
		}
	}

	return tgbotapi.Update{}, false
}

func (b BotReader) StopReceivingUpdates() {
	b.api.StopReceivingUpdates()
}
