package app

import (
	"context"
	"fmt"

	"github.com/DmitySH/tg-gpt/internal/gpt"
	"github.com/DmitySH/tg-gpt/internal/pkg/closer"
	"github.com/DmitySH/tg-gpt/internal/pkg/loggy"
	"github.com/DmitySH/tg-gpt/internal/pkg/secret"
	"github.com/DmitySH/tg-gpt/internal/service"
	"github.com/DmitySH/tg-gpt/internal/store"
	"github.com/DmitySH/tg-gpt/internal/tgbot"
	"github.com/DmitySH/tg-gpt/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	cfg Config
}

func NewApp(config Config) *App {
	return &App{
		cfg: config,
	}
}

func (a *App) Run(_ context.Context) error {
	botAPI, err := initBotAPI()
	if err != nil {
		return err
	}

	chatSessionStorage := store.NewChatSessionStorage(a.cfg.TTL, uint64(a.cfg.MaxCapacity))

	chatGPT, err := gpt.NewOpenAI(gpt.OpenAIConfig{
		APIKey:    secret.GetString("CHAT_GTP_API_KEY"),
		ProxyURL:  secret.GetString("PROXY_URL"),
		ProxyAuth: secret.GetString("PROXY_AUTH"),
	})
	if err != nil {
		return fmt.Errorf("can't init chat gpt: %w", err)
	}

	updateProcessor := service.NewUpdateProcessor(
		service.UpdateProcessorConfig{
			WorkersCount: a.cfg.App.UpdateProcessorWorkerCount,
		},
		usecase.NewOnCommandUsecase(
			chatSessionStorage,
			botAPI,
		),
		usecase.NewOnBasicMessageUsecase(
			chatSessionStorage,
			botAPI,
			chatGPT,
		),
	)

	botReader := tgbot.NewBotReader(tgbot.Config{
		UseWebHook: false,
	}, botAPI, updateProcessor)

	closer.Add(closer.NoErrAdapter(botReader.StopReceivingUpdates))
	closer.Add(closer.NoErrAdapter(updateProcessor.Stop))

	botReader.StartReceivingUpdates()

	closer.Wait()
	loggy.Sync()

	return nil
}

func initBotAPI() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(secret.GetString("BOT_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("can't init bot api: %w", err)
	}

	return bot, nil
}
